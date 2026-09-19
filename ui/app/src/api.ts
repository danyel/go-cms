export type User = {
  id?: string
  name?: string
  email?: string
  picture?: string
}

export type SessionResponse = {
  authenticated?: boolean
  user?: User
  token?: string
  accessToken?: string
  sessionToken?: string
}

const baseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/$/, '')
export const apiConfig = {
  google: import.meta.env.VITE_GOOGLE_AUTH_ENDPOINT ?? '/api/auth/google',
  session: import.meta.env.VITE_SESSION_ENDPOINT ?? '/api/auth/session',
  logout: import.meta.env.VITE_LOGOUT_ENDPOINT ?? '/api/auth/logout',
}

export function apiUrl(path: string) {
  if (/^https?:\/\//i.test(path)) return path
  return `${baseUrl}${path.startsWith('/') ? path : `/${path}`}`
}

export async function requestSession(): Promise<SessionResponse | null> {
  const response = await fetch(apiUrl(apiConfig.session), {
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  if (response.status === 401 || response.status === 404) return null
  if (!response.ok) throw new Error('Unable to check your session.')
  return (await response.json()) as SessionResponse
}

export async function signInWithGoogle(): Promise<SessionResponse> {
  window.location.assign(apiUrl(apiConfig.google))
  return {}
}

export async function signOut() {
  await fetch(apiUrl(apiConfig.logout), {
    method: 'POST',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
}
