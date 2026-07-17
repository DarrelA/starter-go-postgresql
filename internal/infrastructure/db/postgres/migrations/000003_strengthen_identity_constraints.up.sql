ALTER TABLE users ALTER COLUMN user_uuid SET NOT NULL;

ALTER TABLE provider_identities
  ADD CONSTRAINT provider_identities_provider_not_blank CHECK (btrim(provider) <> ''),
  ADD CONSTRAINT provider_identities_subject_not_blank CHECK (btrim(provider_subject) <> '');
