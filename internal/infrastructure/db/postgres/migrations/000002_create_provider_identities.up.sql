ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

CREATE TABLE provider_identities (
  id BIGSERIAL PRIMARY KEY,
  user_uuid UUID NOT NULL REFERENCES users(user_uuid) ON DELETE CASCADE,
  provider VARCHAR(50) NOT NULL,
  provider_subject VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT provider_identities_provider_subject_unique UNIQUE (provider, provider_subject),
  CONSTRAINT provider_identities_provider_user_unique UNIQUE (provider, user_uuid)
);

CREATE INDEX provider_identities_user_uuid_idx ON provider_identities(user_uuid);
