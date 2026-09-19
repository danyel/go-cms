# CMS project

This repository contains a Go backend and React frontend for the CMS.

## Backend

The backend targets Go 1.27.1 and uses a layered design (`internal/web`, `service`, `repository`, `persistence`, and `domain`). PostgreSQL stores application users and sessions; SQLite stores admin records and admin sessions. Goose migrations run automatically at startup.

```sh
go run ./cmd/server
```

For local development, `make up` starts PostgreSQL, `make migrate` applies both database migrations, `make seed` loads repeatable sample content, `make backend` runs the Go server, and `make frontend` runs the React app. `make dev` runs the database, migrations, backend, and frontend together.

Configuration is read from environment variables (see `.env`). The API includes `GET /health`, `GET /api/auth/google`, `GET /api/auth/google/callback`, `GET /api/auth/session`, `POST /api/login`, `GET /api/test/protected`, and the intentionally non-public admin route `POST /admin/login`. Google uses an authorization-code flow when its environment variables are configured; additional providers can implement `identity.IProvider`. The direct login endpoint remains a local integration seam and should not be exposed as a production identity flow.

Content supports one reusable category per document plus any number of badges. The authenticated `GET /api/content` endpoint accepts `category=engineering` and repeated `badge=go&badge=docker` filters; selected badges are combined with AND semantics. `GET /api/content/categories` returns the category dropdown values. The frontend applies these filters immediately when the category changes or when a badge is entered.

The admin route is intentionally not linked from the frontend. Admin credential policy and admin pages are reserved for a later iteration; the current model and service boundary provide the persistence seam.
