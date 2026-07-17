# Deployment and migrations

Schema migrations are embedded in both application binaries. The dedicated `starter-go-postgresql-migrate` command reads only the PostgreSQL environment variables, applies all pending migrations, and exits successfully when the database is current.

## Local development

Development and test environments default to `MIGRATE_ON_STARTUP=true`, preserving the convenient single-process workflow. It can also be set explicitly in local configuration.

Run the migration job independently with:

```sh
make migrate APP_ENV=dev
```

The command uses the same image and PostgreSQL configuration as the application. PostgreSQL advisory locking serializes concurrent migration jobs, so accidentally starting more than one does not apply a migration twice.

## Release ordering

Production defaults to `MIGRATE_ON_STARTUP=false`. A release should follow this order:

1. Publish the immutable application image.
2. Run that image's migration command as a one-shot job.
3. Stop the release if the job exits unsuccessfully.
4. Start or roll out application replicas only after the job succeeds.

The Compose overlay demonstrates this dependency locally:

```sh
cd deployment
APP_ENV=prod docker compose \
  -f docker-compose.yml \
  -f docker-compose.release.yml \
  --profile release up --build app
```

The overlay disables application startup migration and uses `service_completed_successfully`, so the application is not started when the migration job fails.

For an orchestrator, [`migrate-job.example.yml`](../deployment/migrate-job.example.yml) is a minimal Kubernetes Job template. Replace the image and secret names, apply the Job, wait for its `Complete` condition, and only then update the application Deployment. Database credentials should come from the platform's secret store and the Job should use the same release image as the replicas.

Migration execution is transactional. A failed migration rolls back its schema and ledger changes, returns a non-zero process status, and must block the application rollout. Applied migration checksums are immutable; editing an already-applied migration also fails the job.
