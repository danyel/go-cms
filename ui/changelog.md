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

## 2026-09-19 — Content and authorization

### Prompt

Fix the missing protected page and add a blog-like content experience with
infinite scrolling summaries, detail links, capability-controlled editing and
saving. Persist logged-in users so authorities can be assigned. Add a content
data model with creation/edit audit information and a history table recording
changes. Implement the domain model, mappers, web handlers, services,
repositories, migrations, and pages.

### Tasks completed

- Added persisted user role/editor authorization fields and generic `canEdit`
  capability evaluation; administrators can edit through the same capability.
- Added content and content-history domain models, GORM models, explicit
  mappers, repositories, services, and a Goose PostgreSQL migration.
- Added authenticated content listing/detail/update endpoints with pagination,
  validation, 401/403 handling, audit fields, and transactional history
  snapshots.
- Added `/protected` and `/api/protected`; retained `/test` and
  `/api/test/protected` compatibility routes.
- Added the authenticated infinite-scroll content page, detail page, and
  capability-gated edit/save mode.
- Exposed the authenticated user and generic capabilities in the session
  response; no admin link was added to user-facing navigation.
- Revalidated with `go test ./...`, `npm run typecheck`, and `npm run build`.

## 2026-09-19 — Single sign-on

### Prompt

Change the security model: the owner is the only person who enters the application, so use SSO.
When the token is absent the visitor is anonymous; when it is present the owner can change the
application. Remove sign-in entirely, along with the admin database and the setup.

### Tasks completed

- Replaced the Google OAuth flow, direct login, and admin login with a single SSO seam: the
  upstream identity-aware proxy authenticates the owner and the backend verifies the forwarded
  header (`internal/security`).
- Removed the identity package, the auth service, the SQLite admin database, the admin migrations,
  the admin login route, and every Google/OAuth setting and credential file.
- Reads are anonymous; `PUT /api/content/:slug` returns `401` unless the SSO header is present.
- Added the `GET /api/auth/session` capability response (`authenticated`, `canEdit`), which is
  `false`/`false` for anonymous visitors and `true`/`true` for the owner.
- Dropped `users`, `sessions`, and all actor columns from PostgreSQL through migration `004_sso`;
  content history keeps its snapshot timeline without an actor.
- Reworked the React app: no sign-in screen or sign-out button, public content browsing, an
  SSO/anonymous mode badge in the header, and an Edit button gated on `canEdit`.
- Added backend tests for anonymous read-only access and proxy-authenticated editing.
- Revalidated with `gofmt`, `go build ./...`, `go test ./...`, `npm run typecheck`, and `npm run build`.
