/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_SESSION_ENDPOINT?: string
  readonly VITE_CONTENT_ENDPOINT?: string
  /** Google/IdP sign-in URL the header button redirects to; empty hides the button. */
  readonly VITE_SSO_LOGIN_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
