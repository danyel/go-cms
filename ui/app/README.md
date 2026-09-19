# CMS frontend

This Vite/React app lives in `ui/app` and is intentionally independent of the Go backend. Run
`npm install`, then `npm run dev` or `npm run build`.

There is no sign-in screen. Content is public and read-only for anonymous visitors; the Edit
button appears only when the backend reports `canEdit: true`, which happens when the upstream SSO
proxy asserted the owner's identity on the request. `/api/auth/session` returns
`{ "authenticated": bool, "canEdit": bool }` and never creates a session.

The API defaults to same-origin `/api/auth/session` and `/api/content`. Set `VITE_API_BASE_URL`
for a separate backend, or override the endpoint paths with `VITE_SESSION_ENDPOINT` and
`VITE_CONTENT_ENDPOINT`. Requests use cookies (`credentials: include`) so the edge proxy cookie
travels with them.