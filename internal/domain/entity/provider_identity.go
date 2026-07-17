package entity

import "github.com/google/uuid"

const OAuthProviderGoogle = "google"

// ProviderIdentity links a stable external-provider subject to a local user.
type ProviderIdentity struct {
	UserUUID uuid.UUID
	Provider string
	Subject  string
}
