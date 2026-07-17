package repository

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

// ProviderIdentityRepository resolves external identities to local users.
type ProviderIdentityRepository interface {
	FindLinkOrCreate(
		ctx context.Context,
		identity entity.ProviderIdentity,
		candidate *entity.User,
	) (*entity.User, error)
}
