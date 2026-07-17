<a name="readme-top"></a>

- [Intro](#intro)
  - [Overview](#overview)
  - [Documentation](#documentation)
- [Setup](#setup)
  - [Handle Initial Files](#handle-initial-files)
  - [Generate the Private and Public Keys](#generate-the-private-and-public-keys)
    - [Shell Script](#shell-script)
  - [Staticcheck](#staticcheck)
- [Shell](#shell)

# Intro

This is a Go starter application using PostgreSQL for user data and Redis for revocable access and refresh-token sessions. Its package boundaries follow a pragmatic hexagonal architecture: application and domain code define contracts, while HTTP, PostgreSQL, Redis, JWT, and bcrypt are replaceable edge adapters.

## Overview

- **Go and PostgreSQL**: PostgreSQL is the source of truth for users, accessed through a repository contract owned by the core.
- **Rotating token sessions with Redis**: JWTs carry identity claims while Redis atomically rotates single-use refresh tokens, detects replay, and supports family-wide revocation.
- **Unified password and Google login**: Verified Google identities resolve to local users and establish the same cookie-based application sessions as password login.
- **Transport-neutral core**: Application services use `context.Context` and ordinary errors; Fiber-specific request handling and status mapping remain in the HTTP adapter.
- **Explicit composition**: `cmd/auth` constructs infrastructure dependencies, injects them into application services and handlers, and owns startup and graceful shutdown.
- **Testable boundaries**: Small consumer-facing interfaces allow application and HTTP behavior to be tested without live PostgreSQL or Redis instances.

## Documentation

- [`docs`](./docs/README.md) describes the implemented architecture, authentication model, and current limitations.
- [`docs/deployment.md`](./docs/deployment.md) describes one-shot migrations and release ordering.
- HTTP API documentation is maintained as Go documentation beside the transport code. Run `make docs` to open the module documentation locally with the version of pkgsite pinned in `go.mod`.

The first `make docs` run downloads the pinned tool, opens `http://localhost:6060`, and serves documentation until stopped with `Ctrl+C`. Use a different address when needed, for example `make docs DOCS_ADDR=localhost:6061`.

# Setup

## Handle Initial Files

Create ignored local configuration from the sanitized templates:

```sh
make init-env APP_ENV=dev
make init-env APP_ENV=test
```

This creates the respective `.env` and optional pgAdmin server files and generates
local JWT signing keys and a CSRF signing secret. Review the generated files before starting the services,
especially database, pgAdmin, and OAuth2 credentials. Production credentials
and signing keys should come from a secret manager rather than local files.
Configuration is validated at startup. `SEED_DATA=true` loads the development
or test sample users after pending database migrations have been applied; keep
it disabled in production.

## Generate the Private and Public Keys

### Shell Script

```sh
# Rotate the local access/refresh signing keys and CSRF secret.
make rotate-keys APP_ENV=dev

# Format app.log
make lg
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Staticcheck

[Staticcheck](https://staticcheck.dev/) complements `go vet` with additional correctness, performance, simplification, and maintainability checks. Install it as a user-level Go tool, then run it from the repository root:

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck -version
staticcheck ./...
```

Installing the tool does not modify this repository or its `go.mod` file. CI should pin a specific Staticcheck release instead of relying on `@latest`.

# Shell

Project setup, service lifecycle, testing, and maintenance commands are defined in the [`Makefile`](./Makefile).

| Task | Command |
| --- | --- |
| Start the application without pgAdmin | `make up` |
| Start the application with pgAdmin | `make up-pgadmin` |
| Apply pending database migrations | `make migrate APP_ENV=dev` |
| Open local Go documentation | `make docs` |
| Run additional static analysis | `staticcheck ./...` |
| Show the directory tree | `tree -d` |
| Open pgAdmin on macOS | `open http://localhost:5050/browser/` |
| Open a PostgreSQL shell | `cd deployment && APP_ENV=dev docker compose exec postgres psql -U <user> -d <database>` |
| List PostgreSQL databases | `\l` |
| Query users in `psql` | `SELECT * FROM users;` |
| Open a Redis shell | `cd deployment && APP_ENV=dev docker compose exec redis redis-cli` |
| Show Redis server information | `INFO` |
| List session keys | `SCAN 0 MATCH 'session:*'` |
| Inspect a session key | `TYPE <session-key>` then `GET`, `HGETALL`, or `SMEMBERS` as appropriate |

<p align="right">(<a href="#readme-top">back to top</a>)</p>
