# Project documentation

## Guides

- [Architecture](./architecture.md): package boundaries, dependency direction, runtime composition, testing strategy, and current limitations.
- [Authentication](./authentication.md): password login, JWT and Redis session behavior, cookies, and Google OAuth.
- [Deployment](./deployment.md): one-shot migrations, release ordering, and multi-replica rollout guidance.
- API documentation lives beside the implementation as Go documentation. Run `make docs` from the repository root to open the complete module documentation in a browser; the command uses the pkgsite version pinned in `go.mod`.

For local setup, configuration generation, and common commands, see the [root README](../README.md).
