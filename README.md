# CMS project

This repository contains a Go backend and React frontend for the CMS.

## Backend

The backend targets Go 1.27.1 and uses a layered design (`internal/web`, `service`, `repository`, `persistence`, and `domain`). SQLite is the default local database; Goose migrations run automatically at startup.

```sh
go run ./cmd/server
```

Configuration is read from environment variables (see `.env.example`). The API includes `GET /health`, `GET /api/auth/google`, `GET /api/auth/google/callback`, `GET /api/auth/session`, `POST /api/login`, `GET /api/test/protected`, and the intentionally non-public admin route `POST /admin/login`. Google uses an authorization-code flow when its environment variables are configured; additional providers can implement `identity.IProvider`. The direct login endpoint remains a local integration seam and should not be exposed as a production identity flow.

The admin route is intentionally not linked from the frontend. Admin credential policy and admin pages are reserved for a later iteration; the current model and service boundary provide the persistence seam.
