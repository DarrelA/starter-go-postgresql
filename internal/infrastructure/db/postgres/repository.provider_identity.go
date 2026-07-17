package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryUserByProviderIdentity = `
		SELECT u.user_uuid, u.first_name, u.last_name, u.email
		FROM provider_identities i
		JOIN users u ON u.user_uuid = i.user_uuid
		WHERE i.provider = $1 AND i.provider_subject = $2`
	queryUserByOAuthEmail = `
		SELECT user_uuid, first_name, last_name, email
		FROM users WHERE email = $1 FOR UPDATE`
	queryProviderSubjectByUser = `
		SELECT provider_subject FROM provider_identities
		WHERE provider = $1 AND user_uuid = $2`
	queryInsertOAuthUser = `
		INSERT INTO users(first_name, last_name, email, password)
		VALUES ($1, $2, $3, NULL)
		RETURNING user_uuid`
	queryInsertProviderIdentity = `
		INSERT INTO provider_identities(user_uuid, provider, provider_subject)
		VALUES ($1, $2, $3)`
)

// PostgresProviderIdentityRepository resolves OAuth identities transactionally.
type PostgresProviderIdentityRepository struct{ pool *pgxpool.Pool }

func NewProviderIdentityRepository(pool *pgxpool.Pool) repository.ProviderIdentityRepository {
	return &PostgresProviderIdentityRepository{pool: pool}
}

func (r *PostgresProviderIdentityRepository) FindLinkOrCreate(
	ctx context.Context,
	identity entity.ProviderIdentity,
	candidate *entity.User,
) (*entity.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin OAuth identity transaction: %w", err)
	}
	defer tx.Rollback(ctx) // No-op after Commit.

	if user, err := findUserByProviderIdentity(ctx, tx, identity); err == nil {
		return user, nil
	} else if !errors.Is(err, apperror.ErrUserNotFound) {
		return nil, err
	}

	user, err := findUserByOAuthEmail(ctx, tx, candidate.Email)
	switch {
	case err == nil:
		var existingSubject string
		err := tx.QueryRow(ctx, queryProviderSubjectByUser, identity.Provider, *user.UUID).Scan(&existingSubject)
		if err == nil && existingSubject != identity.Subject {
			return nil, apperror.ErrOAuthIdentityConflict
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("check provider identity for user: %w", err)
		}
		identity.UserUUID = *user.UUID
	case errors.Is(err, apperror.ErrUserNotFound):
		user = candidate
		var userUUID uuid.UUID
		if err := tx.QueryRow(
			ctx, queryInsertOAuthUser, user.FirstName, user.LastName, user.Email,
		).Scan(&userUUID); err != nil {
			return nil, mapOAuthPersistenceError("insert OAuth user", err)
		}
		user.UUID = &userUUID
		identity.UserUUID = userUUID
	default:
		return nil, err
	}

	if _, err := tx.Exec(
		ctx, queryInsertProviderIdentity, identity.UserUUID, identity.Provider, identity.Subject,
	); err != nil {
		return nil, mapOAuthPersistenceError("insert provider identity", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit OAuth identity transaction: %w", err)
	}
	return user, nil
}

func findUserByProviderIdentity(
	ctx context.Context,
	tx pgx.Tx,
	identity entity.ProviderIdentity,
) (*entity.User, error) {
	return scanOAuthUser(tx.QueryRow(ctx, queryUserByProviderIdentity, identity.Provider, identity.Subject))
}

func findUserByOAuthEmail(ctx context.Context, tx pgx.Tx, email string) (*entity.User, error) {
	return scanOAuthUser(tx.QueryRow(ctx, queryUserByOAuthEmail, email))
}

type rowScanner interface{ Scan(dest ...any) error }

func scanOAuthUser(row rowScanner) (*entity.User, error) {
	user := &entity.User{}
	if err := row.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, fmt.Errorf("resolve OAuth user: %w", err)
	}
	return user, nil
}

func mapOAuthPersistenceError(operation string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperror.ErrOAuthIdentityConflict
	}
	return fmt.Errorf("%s: %w", operation, err)
}
