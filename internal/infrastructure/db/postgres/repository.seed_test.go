package postgres

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
)

func TestSeedReconcilesEveryUser(t *testing.T) {
	repository := &seedUserRepository{users: map[string]*entity.User{
		"first@example.com": {Email: "first@example.com"},
	}}
	seeder := newTestSeeder(t, `[
		{"first_name":"First","last_name":"User","email":"first@example.com","password":"Password1!"},
		{"first_name":"Second","last_name":"User","email":"second@example.com","password":"Password1!"},
		{"first_name":"Third","last_name":"User","email":"third@example.com","password":"Password1!"}
	]`)

	if err := seeder.Seed(context.Background(), repository); err != nil {
		t.Fatalf("seed users: %v", err)
	}
	if len(repository.users) != 3 {
		t.Fatalf("expected all seed users, got %d", len(repository.users))
	}
	if repository.users["second@example.com"].Password == "Password1!" {
		t.Fatal("expected newly seeded password to be hashed")
	}
}

func TestSeedRecoversAfterPartialFailure(t *testing.T) {
	repository := &seedUserRepository{
		users:       make(map[string]*entity.User),
		failOnEmail: "second@example.com",
	}
	seeder := newTestSeeder(t, `[
		{"first_name":"First","last_name":"User","email":"first@example.com","password":"Password1!"},
		{"first_name":"Second","last_name":"User","email":"second@example.com","password":"Password1!"}
	]`)

	if err := seeder.Seed(context.Background(), repository); err == nil {
		t.Fatal("expected partial seed failure")
	}
	if _, ok := repository.users["first@example.com"]; !ok {
		t.Fatal("expected first user to remain after partial failure")
	}

	repository.failOnEmail = ""
	if err := seeder.Seed(context.Background(), repository); err != nil {
		t.Fatalf("resume seed: %v", err)
	}
	if len(repository.users) != 2 {
		t.Fatalf("expected recovered seed data, got %d users", len(repository.users))
	}
}

func newTestSeeder(t *testing.T, data string) PostgresSeedRepository {
	t.Helper()
	basePath := t.TempDir()
	jsonPath := filepath.Join(basePath, "json")
	if err := os.Mkdir(jsonPath, 0o755); err != nil {
		t.Fatalf("create seed directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(jsonPath, "seed.user.test.json"), []byte(data), 0o600); err != nil {
		t.Fatalf("write seed data: %v", err)
	}
	return PostgresSeedRepository{env: "test", envBasePath: basePath}
}

type seedUserRepository struct {
	users       map[string]*entity.User
	failOnEmail string
}

func (r *seedUserRepository) Save(_ context.Context, user *entity.User) error {
	if user.Email == r.failOnEmail {
		return errors.New("injected save failure")
	}
	if _, exists := r.users[user.Email]; exists {
		return apperror.ErrEmailConflict
	}
	copy := *user
	r.users[user.Email] = &copy
	return nil
}

func (r *seedUserRepository) GetByEmail(_ context.Context, email string) (*entity.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, apperror.ErrUserNotFound
	}
	return user, nil
}

func (r *seedUserRepository) GetByUUID(context.Context, uuid.UUID) (*entity.User, error) {
	return nil, apperror.ErrUserNotFound
}
