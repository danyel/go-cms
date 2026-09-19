export type User = {
  id?: string
  name?: string
  email?: string
  picture?: string
  capabilities?: Record<string, boolean>
}

export type SessionResponse = {
  authenticated?: boolean
  user?: User
  capabilities?: Record<string, boolean>
  canEdit?: boolean
  token?: string
  accessToken?: string
  sessionToken?: string
}

const baseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/$/, '')
export const apiConfig = {
  google: import.meta.env.VITE_GOOGLE_AUTH_ENDPOINT ?? '/api/auth/google',
  session: import.meta.env.VITE_SESSION_ENDPOINT ?? '/api/auth/session',
  logout: import.meta.env.VITE_LOGOUT_ENDPOINT ?? '/api/auth/logout',
  content: import.meta.env.VITE_CONTENT_ENDPOINT ?? '/api/content',
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
  const value = (await response.json()) as SessionResponse & { CanEdit?: boolean; User?: User }
  if (value.User && !value.user) value.user = value.User
  if (value.canEdit === undefined && value.CanEdit !== undefined) value.canEdit = value.CanEdit
  if (value.user) value.user = normalizeUser(value.user)
  return value
}

function normalizeUser(value: User & { ID?: string; Email?: string; Name?: string; Editor?: boolean }): User {
  return { id: value.id ?? value.ID, email: value.email ?? value.Email, name: value.name ?? value.Name, picture: value.picture, capabilities: { ...(value.capabilities ?? {}), ...(value.Editor === true ? { canEdit: true } : {}) } }
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

export type Content = {
  id: string
  title: string
  summary: string
  body: string
  status: string
  createdAt?: string
  updatedAt?: string
}

export type ContentListResponse = {
  items?: Content[]
  content?: Content[]
  data?: Content[]
  nextCursor?: string | null
  cursor?: string | null
  hasMore?: boolean
  offset?: number
  limit?: number
}

function normalizeContent(value: Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string }): Content {
  return { id: value.id ?? value.Slug ?? String(value.ID ?? ''), title: value.title ?? value.Title ?? '', summary: value.summary ?? value.Summary ?? '', body: value.body ?? value.Body ?? '', status: value.status ?? value.Status ?? '', createdAt: value.createdAt, updatedAt: value.updatedAt }
}

async function contentRequest<T>(path = '', init?: RequestInit): Promise<T> {
  const response = await fetch(apiUrl(`${apiConfig.content}${path}`), {
    credentials: 'include',
    ...init,
    headers: { Accept: 'application/json', ...(init?.body ? { 'Content-Type': 'application/json' } : {}), ...init?.headers },
  })
  if (!response.ok) {
    if (response.status === 401) throw new Error('Please sign in to view content.')
    throw new Error((await response.text()) || 'Unable to load content.')
  }
  return (await response.json()) as T
}

export async function listContent(cursor?: string | null): Promise<ContentListResponse> {
  const query = cursor ? `?offset=${encodeURIComponent(cursor)}` : ''
  const response = await contentRequest<ContentListResponse | (Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string })[]>(query)
  if (Array.isArray(response)) return { items: response.map(normalizeContent), nextCursor: null }
  const items = (response.items ?? response.content ?? response.data ?? []).map(normalizeContent)
  const offset = response.offset ?? (cursor ? Number(cursor) : 0)
  const limit = response.limit ?? 20
  return { ...response, items, nextCursor: items.length >= limit ? String(offset + items.length) : null }
}

export function getContent(id: string): Promise<Content> {
  return contentRequest<Content>(`/${encodeURIComponent(id)}`).then(value => normalizeContent(value as Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string }))
}

export type ContentUpdate = Pick<Content, 'title' | 'summary' | 'body' | 'status'>

export function updateContent(id: string, content: ContentUpdate): Promise<Content> {
  return contentRequest<Content>(`/${encodeURIComponent(id)}`, {
    method: 'PUT',
    body: JSON.stringify(content),
  }).then(value => normalizeContent(value as Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string }))
}
