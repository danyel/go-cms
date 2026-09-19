# CMS frontend

This Vite/React app lives in `ui/app` and is intentionally independent of the Go backend. Run
`npm install`, then `npm run dev` or `npm run build`.

The API defaults to same-origin `/api/auth/session`, `/api/auth/google`, and
`/api/auth/logout`. Set `VITE_API_BASE_URL` for a separate backend, or override the endpoint
paths with `VITE_SESSION_ENDPOINT`, `VITE_GOOGLE_AUTH_ENDPOINT`, and `VITE_LOGOUT_ENDPOINT`.
The Google endpoint may return JSON containing `user` and an optional `token`/`accessToken`, or a
redirect response. Sessions use cookies (`credentials: include`) and tokens are retained only for
the current browser tab. Provider-specific login is isolated in `src/api.ts` so Facebook can be
added without changing the routes or UI.
