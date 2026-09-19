export type User = {
  id?: string
  name?: string
  email?: string
  picture?: string
  capabilities?: Record<string, boolean>
}

// The backend never signs users in: an upstream SSO proxy asserts the owner's
// identity. The session endpoint only reports whether that identity is present.
export type SessionResponse = {
  authenticated?: boolean
  canEdit?: boolean
}

const baseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/$/, '')
export const apiConfig = {
  session: import.meta.env.VITE_SESSION_ENDPOINT ?? '/api/auth/session',
  content: import.meta.env.VITE_CONTENT_ENDPOINT ?? '/api/content',
  signIn: import.meta.env.VITE_SSO_LOGIN_URL ?? '',
}

export function apiUrl(path: string) {
  if (/^https?:\/\//i.test(path)) return path
  return `${baseUrl}${path.startsWith('/') ? path : `/${path}`}`
}

// Send the browser to the configured Google/IdP sign-in URL. The identity
// provider authenticates the owner and sends them back here; the app never
// performs the OAuth exchange itself, so no client secret is needed. The
// current page is passed as `rd` for proxies (oauth2-proxy) that honour it and
// ignored by the ones that do not.
export function signInWithGoogle(target = apiConfig.signIn): void {
  if (!target) return
  const returnTo = apiUrl(`${window.location.pathname}${window.location.search}`)
  window.location.assign(`${target}${target.includes('?') ? '&' : '?'}rd=${encodeURIComponent(returnTo)}`)
}

export async function requestSession(): Promise<SessionResponse | null> {
  const response = await fetch(apiUrl(apiConfig.session), {
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  if (response.status === 404) return null
  if (!response.ok) throw new Error('Unable to check your session.')
  const value = (await response.json()) as SessionResponse & { CanEdit?: boolean }
  if (value.canEdit === undefined && value.CanEdit !== undefined) value.canEdit = value.CanEdit
  return value
}

export type Content = {
  id: string
  title: string
  summary: string
  body: string
  status: string
  category?: string
  badges?: string[]
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

function normalizeContent(value: Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string; Category?: string; Badges?: string[] }): Content {
  return { id: value.id ?? value.Slug ?? String(value.ID ?? ''), title: value.title ?? value.Title ?? '', summary: value.summary ?? value.Summary ?? '', body: value.body ?? value.Body ?? '', status: value.status ?? value.Status ?? '', category: value.category ?? value.Category, badges: value.badges ?? value.Badges ?? [], createdAt: value.createdAt, updatedAt: value.updatedAt }
}

async function contentRequest<T>(path = '', init?: RequestInit): Promise<T> {
  const response = await fetch(apiUrl(`${apiConfig.content}${path}`), {
    credentials: 'include',
    ...init,
    headers: { Accept: 'application/json', ...(init?.body ? { 'Content-Type': 'application/json' } : {}), ...init?.headers },
  })
  if (!response.ok) {
    if (response.status === 401) throw new Error('Your SSO session is required to edit content.')
    throw new Error((await response.text()) || 'Unable to load content.')
  }
  return (await response.json()) as T
}

export async function listContent(cursor?: string | null, category = '', status = '', badges: string[] = []): Promise<ContentListResponse> {
  const params = new URLSearchParams()
  if (cursor) params.set('offset', cursor)
  if (category) params.set('category', category)
  if (status) params.set('status', status)
  badges.forEach(badge => params.append('badge', badge))
  const query = params.toString() ? `?${params.toString()}` : ''
  const response = await contentRequest<ContentListResponse | (Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string })[]>(query)
  if (Array.isArray(response)) return { items: response.map(normalizeContent), nextCursor: null }
  const items = (response.items ?? response.content ?? response.data ?? []).map(normalizeContent)
  const offset = response.offset ?? (cursor ? Number(cursor) : 0)
  const limit = response.limit ?? 20
  return { ...response, items, nextCursor: items.length >= limit ? String(offset + items.length) : null }
}

export type Category = { id: string; slug: string; name: string }
export type Badge = { id: string; name: string }

export async function listCategories(): Promise<Category[]> {
  const response = await contentRequest<{ items?: Array<Category & { ID?: string; Slug?: string; Name?: string }> }>('/categories')
  return (response.items ?? []).map(item => ({ id: item.id ?? item.ID ?? '', slug: item.slug ?? item.Slug ?? '', name: item.name ?? item.Name ?? '' }))
}

export async function listBadges(): Promise<Badge[]> {
  const response = await contentRequest<{ items?: Array<Badge & { ID?: string; Name?: string }> }>('/badges')
  return (response.items ?? []).map(item => ({ id: item.id ?? item.ID ?? '', name: item.name ?? item.Name ?? '' }))
}

export function getContent(id: string): Promise<Content> {
  return contentRequest<Content>(`/${encodeURIComponent(id)}`).then(value => normalizeContent(value as Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string }))
}

export type ContentUpdate = Pick<Content, 'title' | 'summary' | 'body' | 'status' | 'badges'>

export function updateContent(id: string, content: ContentUpdate): Promise<Content> {
  return contentRequest<Content>(`/${encodeURIComponent(id)}`, {
    method: 'PUT',
    body: JSON.stringify(content),
  }).then(value => normalizeContent(value as Content & { ID?: string; Slug?: string; Title?: string; Summary?: string; Body?: string; Status?: string }))
}