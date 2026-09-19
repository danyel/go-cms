# CMS frontend

This Vite/React app lives in `ui/app` and is intentionally independent of the Go backend. Run
`npm install`, then `npm run dev` or `npm run build`.

There is no sign-in screen and no session is ever created. Content is public and read-only for
anonymous visitors; the Edit button appears only when the backend reports `canEdit: true`, which
happens when the upstream SSO proxy asserted the owner's identity on the request.
`/api/auth/session` returns `{ "authenticated": bool, "canEdit": bool }`.

Anonymous visitors also see a **Sign in with Google** button in the header. It redirects the
browser to `VITE_SSO_LOGIN_URL` (with the current page passed as `rd`), where the identity
provider authenticates you and sends you back with the identity header set. Leave
`VITE_SSO_LOGIN_URL` empty to hide the button — useful when the proxy gates the whole app
automatically instead. The app never performs an OAuth exchange, so it holds no client secret.

The API defaults to same-origin `/api/auth/session` and `/api/content`. Set `VITE_API_BASE_URL`
for a separate backend, or override the endpoint paths with `VITE_SESSION_ENDPOINT` and
`VITE_CONTENT_ENDPOINT`. Requests use cookies (`credentials: include`) so the edge proxy cookie
travels with them.