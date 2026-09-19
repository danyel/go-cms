import { useEffect, useRef, useState } from 'react'
import { Link, Navigate, Route, Routes, useLocation, useNavigate, useParams } from 'react-router-dom'
import { getContent, listBadges, listCategories, listContent, updateContent, type Badge, type Category, type Content, type ContentUpdate } from './api'
import { useAuth } from './auth'

function Shell({ children }: { children: React.ReactNode }) {
  const { user, signOut } = useAuth()
  return <main className="shell">
    <header className="header"><Link className="brand" to="/">CMS</Link>{user && <button className="text-button" onClick={() => void signOut()}>Sign out</button>}</header>
    {children}
  </main>
}

function Login() {
  const { user, loading, error, signIn } = useAuth()
  if (user) return <Navigate to="/protected" replace />
  if (loading) return <Shell><div className="card"><p>Checking your session…</p></div></Shell>
  return <Shell><section className="card auth-card">
    <span className="eyebrow">Welcome</span><h1>Sign in to CMS</h1>
    <p className="muted">Sign in securely with your external identity provider.</p>
    {error && <p className="error" role="alert">{error}</p>}
    <button className="google-button" onClick={() => void signIn()}><span aria-hidden="true">G</span> Continue with Google</button>
  </section></Shell>
}

function RequireAuth({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth()
  const location = useLocation()
  if (loading) return <Shell><div className="card"><p>Loading…</p></div></Shell>
  if (!user) return <Navigate to="/" state={{ from: location }} replace />
  return <>{children}</>
}

function ContentCard({ item }: { item: Content }) {
  return <article className="content-card">
    <span className="status">{item.status}</span>
    <h2><Link to={`/content/${item.id}`}>{item.title}</Link></h2>
    <p>{item.summary}</p>
    {item.badges && item.badges.length > 0 && <div className="badge-list">{item.badges.map(badge => <span className="badge" key={badge}>{badge}</span>)}</div>}
    <Link className="read-link" to={`/content/${item.id}`}>Read more <span aria-hidden="true">→</span></Link>
  </article>
}

function ContentList() {
  const [items, setItems] = useState<Content[]>([])
  const [cursor, setCursor] = useState<string | null | undefined>(undefined)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [badgeOptions, setBadgeOptions] = useState<Badge[]>([])
  const [category, setCategory] = useState('')
  const [badges, setBadges] = useState<string[]>([])
  const [badgeInput, setBadgeInput] = useState('')
  const loadingMore = useRef(false)
  const requestVersion = useRef(0)
  const sentinel = useRef<HTMLDivElement>(null)

  const load = async (next = false) => {
    if (loadingMore.current || (next && !cursor)) return
    const version = next ? requestVersion.current : ++requestVersion.current
    loadingMore.current = true
    setLoading(true)
    try {
      const response = await listContent(next ? cursor : undefined, category, badges)
      if (version !== requestVersion.current) return
      const page = Array.isArray(response) ? response : response.items ?? response.content ?? response.data ?? []
      setItems(current => next ? [...current, ...page] : page)
      setCursor(Array.isArray(response) ? null : response.nextCursor ?? response.cursor ?? (response.hasMore ? undefined : null))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unable to load content.')
    } finally {
      if (version === requestVersion.current) {
        loadingMore.current = false
        setLoading(false)
      }
    }
  }
  useEffect(() => {
    void Promise.all([listCategories(), listBadges()])
      .then(([nextCategories, nextBadges]) => { setCategories(nextCategories); setBadgeOptions(nextBadges) })
      .catch(() => { setCategories([]); setBadgeOptions([]) })
  }, [])
  useEffect(() => {
    loadingMore.current = false
    setItems([])
    setCursor(undefined)
    void load()
  }, [category, badges.join(',')])
  useEffect(() => {
    const node = sentinel.current
    if (!node) return
    const observer = new IntersectionObserver(entries => { if (entries[0]?.isIntersecting) void load(true) }, { rootMargin: '240px' })
    observer.observe(node)
    return () => observer.disconnect()
  }, [cursor])

  return <Shell><section className="content-page">
    <div className="page-heading"><div><span className="eyebrow">Journal</span><h1>Latest content</h1></div></div>
    <div className="filters" aria-label="Content filters">
      <label className="category-filter">Category
        <select value={category} onChange={event => setCategory(event.target.value)}>
          <option value="">All categories</option>
          {categories.map(item => <option value={item.slug} key={item.slug}>{item.name}</option>)}
        </select>
      </label>
      <label className="badge-filter">Badges
        <div className="badge-input">
          {badges.map(badge => <span className="filter-pill" key={badge}>{badge}<button type="button" aria-label={`Remove ${badge}`} onClick={() => setBadges(current => current.filter(value => value !== badge))}>×</button></span>)}
          <input list="badge-options" value={badgeInput} placeholder="Type a badge and press Enter" onChange={event => setBadgeInput(event.target.value)} onKeyDown={event => {
            if (event.key !== 'Enter') return
            event.preventDefault()
            const badge = badgeInput.trim().toLowerCase()
            if (badge && !badges.includes(badge)) setBadges(current => [...current, badge])
            setBadgeInput('')
          }} />
          <datalist id="badge-options">{badgeOptions.filter(item => !badges.includes(item.name)).map(item => <option value={item.name} key={item.id}>{item.name}</option>)}</datalist>
        </div>
      </label>
    </div>
    {error && <p className="error" role="alert">{error}</p>}
    {!loading && !error && items.length === 0 && <div className="card empty-state"><h2>No content yet</h2><p className="muted">There is nothing published here yet.</p></div>}
    <div className="content-grid">{items.map(item => <ContentCard item={item} key={item.id} />)}</div>
    <div ref={sentinel} className="list-status" aria-live="polite">{loading ? <p className="muted">Loading content…</p> : cursor ? <p className="muted">Loading more…</p> : items.length > 0 ? <p className="muted">You’ve reached the end.</p> : null}</div>
  </section></Shell>
}

function ContentDetail() {
  const { id = '' } = useParams()
  const { user } = useAuth()
  const canEdit = Boolean(user?.capabilities?.canEdit)
  const [item, setItem] = useState<Content | null>(null)
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState<ContentUpdate | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()
  useEffect(() => {
    let active = true
    setLoading(true)
    void getContent(id).then(value => { if (active) { setItem(value); setDraft({ title: value.title, summary: value.summary, body: value.body, status: value.status }) } })
      .catch(err => active && setError(err instanceof Error ? err.message : 'Unable to load content.'))
      .finally(() => active && setLoading(false))
    return () => { active = false }
  }, [id])
  if (loading) return <Shell><div className="card"><p>Loading content…</p></div></Shell>
  if (error || !item || !draft) return <Shell><div className="card"><p className="error" role="alert">{error ?? 'Content not found.'}</p><Link to="/content">Back to content</Link></div></Shell>
  const save = async () => {
    setSaving(true); setError(null)
    try { const saved = await updateContent(id, draft); setItem(saved); setDraft({ title: saved.title, summary: saved.summary, body: saved.body, status: saved.status }); setEditing(false) }
    catch (err) { setError(err instanceof Error ? err.message : 'Unable to save changes.') }
    finally { setSaving(false) }
  }
  return <Shell><article className="detail-page">
    <Link className="back-link" to="/content">← All content</Link>
    {error && <p className="error" role="alert">{error}</p>}
    {editing ? <div className="editor">
      <label>Title<input value={draft.title} onChange={e => setDraft({ ...draft, title: e.target.value })} /></label>
      <label>Summary<textarea rows={3} value={draft.summary} onChange={e => setDraft({ ...draft, summary: e.target.value })} /></label>
      <label>Body<textarea rows={12} value={draft.body} onChange={e => setDraft({ ...draft, body: e.target.value })} /></label>
      <label>Status<select value={draft.status} onChange={e => setDraft({ ...draft, status: e.target.value })}><option value="draft">Draft</option><option value="published">Published</option><option value="archived">Archived</option></select></label>
      <div className="actions"><button className="primary-button" disabled={saving} onClick={() => void save()}>{saving ? 'Saving…' : 'Save changes'}</button><button className="secondary-button" onClick={() => setEditing(false)}>Cancel</button></div>
    </div> : <><div className="detail-heading"><span className="status">{item.status}</span><h1>{item.title}</h1><p className="muted">{item.summary}</p>{canEdit && <button className="primary-button edit-button" onClick={() => setEditing(true)}>Edit</button>}</div><div className="body-copy">{item.body}</div></>}
  </article></Shell>
}

function ProtectedTest() {
  const { user } = useAuth()
  return <Shell><section className="card protected-card"><span className="eyebrow">Authenticated</span><h1>Protected content</h1><p>You are signed in and can access protected content.</p><div className="profile">{user?.email ?? user?.name ? `Signed in as ${user.email ?? user.name}` : 'Your session is active.'}</div><Link className="primary-button inline-button" to="/content">Browse content</Link></section></Shell>
}

export function App() {
  return <Routes>
    <Route path="/" element={<Login />} />
    <Route path="/protected" element={<RequireAuth><ProtectedTest /></RequireAuth>} />
    <Route path="/test" element={<RequireAuth><ProtectedTest /></RequireAuth>} />
    <Route path="/content" element={<RequireAuth><ContentList /></RequireAuth>} />
    <Route path="/content/:id" element={<RequireAuth><ContentDetail /></RequireAuth>} />
    <Route path="*" element={<Navigate to="/" replace />} />
  </Routes>
}
