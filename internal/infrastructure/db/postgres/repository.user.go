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

// PostgresUserRepository persists users in PostgreSQL.
type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

// NewUserRepository constructs a PostgreSQL user repository.
func NewUserRepository(dbpool *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{dbpool}
}

const (
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING user_uuid;"
	queryGetUser     = "SELECT user_uuid, first_name, last_name, email, COALESCE(password, '') FROM users WHERE email=$1;"
	queryGetUserByID = "SELECT user_uuid, first_name, last_name, email FROM users WHERE user_uuid=$1;"
)

func (ur PostgresUserRepository) Save(ctx context.Context, user *entity.User) error {
	var lastInsertUUID uuid.UUID
	err := ur.dbpool.QueryRow(ctx, queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertUUID)

	if err != nil {
		var pgErr *pgconn.PgError

		// Check if `err` can be cast to `*pgconn.PgError` using `errors.As` before
		// attempting to access any fields of `pgErr`
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return apperror.ErrEmailConflict
			}
		}
		return fmt.Errorf("insert user: %w", err)
	}

	user.UUID = &lastInsertUUID
	return nil
}

func (ur PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := ur.dbpool.QueryRow(ctx, queryGetUser, email).
		Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (ur PostgresUserRepository) GetByUUID(ctx context.Context, userUUID uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	result := ur.dbpool.QueryRow(ctx, queryGetUserByID, userUUID)
	if err := result.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by UUID: %w", err)
	}

	return user, nil
}
