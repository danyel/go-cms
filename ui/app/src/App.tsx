import {useEffect, useRef, useState} from 'react'
import {Link, Navigate, Route, Routes, useParams} from 'react-router-dom'
import {
    type Badge,
    type Category,
    type Content,
    type ContentUpdate,
    getContent,
    listDemoUsers,
    listBadges,
    listCategories,
    listContent,
    updateContent
} from './api'
import {useAuth} from './auth'

function Shell({children}: { children: React.ReactNode }) {
    const {authenticated, loading, canSignIn, signIn, username, sudo, signOut} = useAuth()
    const [sudoOpen, setSudoOpen] = useState(false)
    const [helpOpen, setHelpOpen] = useState(false)
    const [users, setUsers] = useState<{username: string; name: string}[]>([])
    const [user, setUser] = useState('')
    const [password, setPassword] = useState('')
    const [sudoError, setSudoError] = useState('')
    useEffect(() => { void listDemoUsers().then(setUsers) }, [])
    const [dark, setDark] = useState(() => window.localStorage.getItem('urpi-theme') !== 'light')
    useEffect(() => {
        document.documentElement.dataset.theme = dark ? 'dark' : 'light'
        window.localStorage.setItem('urpi-theme', dark ? 'dark' : 'light')
    }, [dark])
    return <main className="shell">
        <header className="header"><Link className="brand" to="/content"><span className="brand-prompt">~/</span> Urpi's
            backlog<span className="brand-cursor">_</span><small>linux · go · java · ideas</small></Link>
            <div className="header-actions">{!loading &&
                <span className="mode-badge">{authenticated ? `${username ?? 'Owner'} · sudo` : 'Anonymous · read-only'}</span>}
                <button className="theme-button" type="button" aria-label="Demo users" onClick={() => setHelpOpen(value => !value)}>?</button>
                {!loading && !authenticated && users.length > 0 &&
                    <button className="theme-button" type="button" onClick={() => setSudoOpen(true)}>sudo</button>}
                {!loading && authenticated && users.length > 0 &&
                    <button className="theme-button" type="button" onClick={signOut}>logout</button>}
                {!loading && !authenticated && canSignIn &&
                    <button className="google-button session-button" type="button" onClick={signIn}><span
                        aria-hidden="true">G</span> Sign in with Google</button>}
                <button className="theme-button" type="button" aria-label={`Switch to ${dark ? 'light' : 'dark'} theme`}
                        onClick={() => setDark(value => !value)}>{dark ? '☀ Light' : '☾ Dark'}</button>
            </div>
        </header>
        {helpOpen && users.length > 0 && <aside className="card terminal-panel">
            <strong>sudo users</strong><p className="muted">Demo only: each password is the same as the username.</p>
            {users.map(item => <div key={item.username}><code>sudo {item.username}</code> · {item.name}</div>)}
        </aside>}
        {sudoOpen && <div className="terminal-backdrop" role="presentation" onMouseDown={event => {
            if (event.target === event.currentTarget) setSudoOpen(false)
        }}>
            <form className="terminal-window" role="dialog" aria-modal="true" aria-labelledby="sudo-title" onSubmit={event => {
                event.preventDefault()
                setSudoError('')
                void sudo(user, password).then(() => { setSudoOpen(false); setPassword('') }).catch(error => setSudoError(error instanceof Error ? error.message : 'sudo failed'))
            }}>
                <div className="terminal-titlebar">
                    <span id="sudo-title">urpi@backlog: ~</span>
                    <button type="button" className="terminal-close" aria-label="Close sudo terminal"
                            onClick={() => setSudoOpen(false)}>×</button>
                </div>
                <div className="terminal-body">
                    <p><span className="terminal-prompt">$</span> sudo login</p>
                    <label><span className="terminal-prompt">$</span> sudo <input autoFocus value={user}
                        onChange={event => setUser(event.target.value)} placeholder="user" aria-label="Username"/></label>
                    <label><span className="terminal-prompt">$</span> password: <input type="password" value={password}
                        onChange={event => setPassword(event.target.value)} aria-label="Password"/></label>
                    {sudoError && <p className="terminal-error">{sudoError}</p>}
                    <button className="terminal-submit" type="submit"><span className="terminal-prompt">$</span> Enter</button>
                </div>
            </form>
        </div>}
        {children}
    </main>
}

function ContentCard({item}: { item: Content }) {
    return <article className="content-card">
        <span className="status">{item.status}</span>
        <h2><Link to={`/content/${item.id}`}>{item.title}</Link></h2>
        <p>{item.summary}</p>
        <p className="content-meta">
            by {item.author ?? 'unknown'}
            {item.publishedAt && <> · published {new Date(item.publishedAt).toLocaleDateString()}</>}
        </p>
        {item.badges && item.badges.length > 0 &&
            <div className="badge-list">{item.badges.map(badge => <span className="badge"
                                                                        key={badge}>{badge}</span>)}</div>}
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
    const [status, setStatus] = useState('')
    const [badges, setBadges] = useState<string[]>([])
    const [badgeInput, setBadgeInput] = useState('')
    const loadingMore = useRef(false)
    const requestVersion = useRef(0)
    const sentinel = useRef<HTMLDivElement>(null)
    const filters = {category, status, badges}

    const load = async (next = false, activeFilters = filters) => {
        if (loadingMore.current || (next && !cursor)) return
        const version = next ? requestVersion.current : ++requestVersion.current
        loadingMore.current = true
        setLoading(true)
        try {
            const response = await listContent(next ? cursor : undefined, activeFilters.category, activeFilters.status, activeFilters.badges)
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
            .then(([nextCategories, nextBadges]) => {
                setCategories(nextCategories);
                setBadgeOptions(nextBadges)
            })
            .catch(() => {
                setCategories([]);
                setBadgeOptions([])
            })
    }, [])
    useEffect(() => {
        loadingMore.current = false
        setItems([])
        setCursor(undefined)
        void load(false, {category, status, badges})
    }, [category, status, badges.join(',')])
    useEffect(() => {
        const node = sentinel.current
        if (!node) return
        const observer = new IntersectionObserver(entries => {
            if (entries[0]?.isIntersecting) void load(true)
        }, {rootMargin: '240px'})
        observer.observe(node)
        return () => observer.disconnect()
    }, [cursor, category, status, badges.join(',')])
    return <Shell>
        <section className="content-page">
            <div className="page-heading">
                <div><span className="eyebrow">$ cd ~/backlog && ls</span></div>
            </div>
            <div className="filters" aria-label="Content filters">
                <label className="status-filter">Status
                    <select value={status} onChange={event => setStatus(event.target.value)}>
                        <option value="">All statuses</option>
                        <option value="published">Published</option>
                        <option value="draft">Draft</option>
                        <option value="review">Review</option>
                        <option value="archived">Archived</option>
                    </select>
                </label>
                <label className="category-filter">Category
                    <select value={category} onChange={event => setCategory(event.target.value)}>
                        <option value="">All categories</option>
                        {categories.map(item => <option value={item.slug} key={item.slug}>{item.name}</option>)}
                    </select>
                </label>
                <label className="badge-filter">Badges
                    <div className="badge-input">
                        {badges.map(badge => <span className="filter-pill" key={badge}>{badge}
                            <button type="button" aria-label={`Remove ${badge}`}
                                    onClick={() => setBadges(current => current.filter(value => value !== badge))}>×</button></span>)}
                        <input list="badge-options" value={badgeInput} placeholder="Type a badge and press Enter"
                               onChange={event => setBadgeInput(event.target.value)} onKeyDown={event => {
                            if (event.key !== 'Enter') return
                            event.preventDefault()
                            const badge = badgeInput.trim().toLowerCase()
                            if (badge && !badges.includes(badge)) setBadges(current => [...current, badge])
                            setBadgeInput('')
                        }}/>
                        <datalist
                            id="badge-options">{badgeOptions.filter(item => !badges.includes(item.name)).map(item =>
                            <option value={item.name} key={item.id}>{item.name}</option>)}</datalist>
                    </div>
                </label>
            </div>
            {error && <p className="error" role="alert">{error}</p>}
            {!loading && !error && items.length === 0 &&
                <div className="card empty-state"><h2>No content yet</h2><p className="muted">There is nothing published
                    here yet.</p></div>}
            <div className="content-grid">{items.map(item => <ContentCard item={item} key={item.id}/>)}</div>
            <div ref={sentinel} className="list-status" aria-live="polite">{loading ?
                <p className="muted">Loading content…</p> : cursor ?
                    <p className="muted">Loading more…</p> : items.length > 0 ?
                        <p className="muted">You’ve reached the end.</p> : null}</div>
        </section>
    </Shell>
}

function ContentDetail() {
    const {id = ''} = useParams()
    const {canEdit, username} = useAuth()
    const [item, setItem] = useState<Content | null>(null)
    const [editing, setEditing] = useState(false)
    const [draft, setDraft] = useState<ContentUpdate | null>(null)
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [badgeOptions, setBadgeOptions] = useState<Badge[]>([])
    const [badgeInput, setBadgeInput] = useState('')
    useEffect(() => {
        let active = true
        setLoading(true)
        void Promise.all([getContent(id), listBadges()]).then(([value, options]) => {
            if (active) {
                setItem(value)
                setBadgeOptions(options)
                setDraft({
                    title: value.title,
                    summary: value.summary,
                    body: value.body,
                    status: value.status,
                    badges: value.badges ?? []
                })
            }
        })
            .catch(err => active && setError(err instanceof Error ? err.message : 'Unable to load content.'))
            .finally(() => active && setLoading(false))
        return () => {
            active = false
        }
    }, [id])
    if (loading) return <Shell>
        <div className="card"><p>Loading content…</p></div>
    </Shell>
    if (error || !item || !draft) return <Shell>
        <div className="card"><p className="error" role="alert">{error ?? 'Content not found.'}</p><Link to="/content">Back
            to content</Link></div>
    </Shell>
    const save = async () => {
        setSaving(true);
        setError(null)
        try {
            const saved = await updateContent(id, draft);
            setItem(saved);
            setDraft({
                title: saved.title,
                summary: saved.summary,
                body: saved.body,
                status: saved.status,
                badges: saved.badges ?? []
            });
            setEditing(false)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Unable to save changes.')
        } finally {
            setSaving(false)
        }
    }
    return <Shell>
        <article className="detail-page">
            <Link className="back-link" to="/content">← All content</Link>
            {error && <p className="error" role="alert">{error}</p>}
            {editing ? <div className="editor">
                <label>Title<input value={draft.title}
                                   onChange={e => setDraft({...draft, title: e.target.value})}/></label>
                <label>Summary<textarea rows={3} value={draft.summary}
                                        onChange={e => setDraft({...draft, summary: e.target.value})}/></label>
                <label>Body<textarea rows={12} value={draft.body}
                                     onChange={e => setDraft({...draft, body: e.target.value})}/></label>
                <label>Status<select value={draft.status} onChange={e => setDraft({...draft, status: e.target.value})}>
                    <option value="draft">Draft</option>
                    <option value="published">Published</option>
                    <option value="archived">Archived</option>
                </select></label>
                <label>Badges
                    <div className="badge-input">
                        {(draft.badges ?? []).map(badge => <span className="filter-pill" key={badge}>{badge}
                            <button type="button" aria-label={`Remove ${badge}`} onClick={() => setDraft({
                                ...draft,
                                badges: (draft.badges ?? []).filter(value => value !== badge)
                            })}>×</button></span>)}
                        <input list="detail-badge-options" value={badgeInput} placeholder="Type a badge and press Enter"
                               onChange={e => setBadgeInput(e.target.value)} onKeyDown={e => {
                            if (e.key !== 'Enter') return
                            e.preventDefault()
                            const badge = badgeInput.trim().toLowerCase()
                            if (badge && !(draft.badges ?? []).includes(badge)) setDraft({
                                ...draft,
                                badges: [...(draft.badges ?? []), badge]
                            })
                            setBadgeInput('')
                        }}/>
                        <datalist
                            id="detail-badge-options">{badgeOptions.filter(option => !(draft.badges ?? []).includes(option.name)).map(option =>
                            <option value={option.name} key={option.id}/>)}</datalist>
                    </div>
                </label>
                <div className="actions">
                    <button className="primary-button" disabled={saving}
                            onClick={() => void save()}>{saving ? 'Saving…' : 'Save changes'}</button>
                    <button className="secondary-button" onClick={() => setEditing(false)}>Cancel</button>
                </div>
            </div> : <>
                <div className="detail-heading"><span className="status">{item.status}</span><h1>{item.title}</h1><p
                    className="muted">{item.summary}</p><p className="muted">by {item.author ?? 'unknown'}{item.publishedAt ? ` · published ${new Date(item.publishedAt).toLocaleDateString()}` : ''}</p>{item.badges && item.badges.length > 0 &&
                    <div className="badge-list">{item.badges.map(badge => <span className="badge"
                                                                                key={badge}>{badge}</span>)}</div>}{canEdit && username === item.author &&
                    <button className="primary-button edit-button" onClick={() => setEditing(true)}>Edit</button>}</div>
                <div className="body-copy">{item.body}</div>
            </>}
        </article>
    </Shell>
}

function ProtectedTest() {
    return <Shell>
        <section className="card protected-card"><span className="eyebrow">Owner session</span><h1>Protected
            content</h1><p>Your SSO proxy asserted the owner identity for this request.</p><Link
            className="primary-button inline-button" to="/content">Browse content</Link></section>
    </Shell>
}

export function App() {
    return <Routes>
        <Route path="/" element={<Navigate to="/content" replace/>}/>
        <Route path="/protected" element={<ProtectedTest/>}/>
        <Route path="/test" element={<ProtectedTest/>}/>
        <Route path="/content" element={<ContentList/>}/>
        <Route path="/content/:id" element={<ContentDetail/>}/>
        <Route path="*" element={<Navigate to="/content" replace/>}/>
    </Routes>
}
