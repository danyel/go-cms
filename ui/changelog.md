# Prompt and implementation changelog

## 2026-09-19

### Prompt

Create a CMS React application with a Go backend, initialize a Git repository, and
prepare the security foundation first. Add an extensible external identity-provider
login beginning with Google, a protected test page, and an admin login at
`/admin/login` whose link is never exposed in user-facing pages. Use Goose for
migrations and GORM for persistence. Structure the backend into web handlers,
web-to-domain mapping, domain-only services, database models, repositories, and
domain/database mappers. Prefix interfaces with `I`. Use Go 1.27.1 and the latest
React. Keep this changelog in `ui/changelog.md`.

### Tasks completed

- Initialized a local Git repository on the `main` branch; GitHub remote
  configuration is intentionally left for later.
- Added a Go 1.27.1 backend module with configuration via environment variables.
- Added layered backend packages for domain, web handlers and mapping, services,
  repositories, persistence models, and persistence mapping.
- Added GORM with SQLite as the local database driver and Goose migrations for
  users, administrators, and sessions.
- Added an OAuth provider abstraction and Google OAuth authorization-code flow with
  state validation, userinfo lookup, and an HTTP-only session cookie.
- Added health, session, Google auth, protected test, and hidden admin login routes.
- Added a React/Vite frontend with a public Google sign-in page, protected `/test`
  page, and `/admin/login` route without a visible navigation link.
- Added environment examples, ignore rules, README documentation, and backend
  unit tests.
- Verified Go tests and prepared frontend typecheck/build scripts.
