/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_GOOGLE_AUTH_ENDPOINT?: string
  readonly VITE_SESSION_ENDPOINT?: string
  readonly VITE_LOGOUT_ENDPOINT?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
