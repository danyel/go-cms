# CMS project

This repository contains a Go backend and React frontend for the CMS.

## Backend

The backend targets Go 1.27.1 and uses a layered design (`internal/web`,
`service`, and `security`). Content is seeded in memory and can be mirrored to
a JSON file configured through `CMS_CONTENT_FILE`.

```sh
go run ./cmd/server
```

Configuration is read from environment variables (see `.env`). The API includes `GET /health`, `GET /api/auth/session`, `GET /api/protected`, `GET /api/content`, `GET /api/content/categories`, `GET /api/content/badges`, and `GET|PUT /api/content/:slug`.

## Single sign-on

The CMS never signs anyone in and stores no users, admins, or sessions. An identity-aware proxy (oauth2-proxy, Authelia, Cloudflare Access, ...) authenticates the owner and forwards a header on every accepted request:

- **Header present** → the request is the owner and can edit content (`canEdit: true`).
- **Header absent** → the request is anonymous and content is read-only (`canEdit: false`).

Configure the header with `CMS_SSO_HEADER` (default `X-SSO-Token`). Set `CMS_SSO_TOKEN` to a shared secret to require an exact value; leave it empty to trust any non-empty forwarded value. `CMS_SSO_COOKIE` optionally names a cookie fallback.

Because the backend trusts the forwarded header, it must only be reachable through the proxy. Bind it to the loopback interface or a private network, and make sure the proxy strips any client-supplied copy of that header.

Anonymous visitors get a **Sign in with Google** button in the header. It redirects the browser to `VITE_SSO_LOGIN_URL` (frontend setting), where the provider authenticates the owner and returns to the page; the CMS itself never performs an OAuth exchange and holds no client secret. Leave the setting empty to hide the button when the proxy gates the whole app automatically.

## Content

Content supports one reusable category per document plus any number of badges. The anonymous `GET /api/content` endpoint accepts `category=engineering` and repeated `badge=go&badge=docker` filters; selected badges are combined with AND semantics. `GET /api/content/categories` returns the category dropdown values. The frontend applies these filters immediately when the category changes or when a badge is entered.

Editing requires the SSO identity: `PUT /api/content/:slug` returns `401` without it. `GET /api/content/:slug` and the taxonomy endpoints are public.

```sh
make backend   # run the Go server
make frontend  # run the React app
make run       # build the frontend, then run the server
```

## Rancher deployment

The Helm chart at `deploy/helm/cms` provides two profiles:

- `values-development.yaml` runs the demo profile on NodePort `30081` with
  ephemeral content.
- `values-production.yaml` runs the seeded demo application on NodePort
  `31373` with a persistent content volume. Host Nginx proxies `cms.urpi.be`
  to this port.

Validate or deploy either profile with:

```sh
make helm-lint
make helm-deploy-development
make helm-deploy-production
```

The deploy workflow automatically deploys a successful `main` image build to
production and also supports manual development or production deployments.
Create GitHub environments named `development` and `production`. Each needs
the secrets `RANCHER_KUBE_CONFIG_BASE64`, `REGISTRY_USERNAME`, and
`REGISTRY_PASSWORD`. `CMS_SSO_TOKEN` is optional. Environment variables
`RANCHER_NAMESPACE`, `CMS_ALLOWED_ORIGINS`, `CMS_SSO_HEADER`, and
`CMS_SSO_COOKIE` can override the chart defaults.