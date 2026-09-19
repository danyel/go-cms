import {useEffect, useRef, useState} from 'react'
import {Link, Navigate, Route, Routes, useParams} from 'react-router-dom'
import {
    type Content,
    type ContentUpdate,
    getContent,
    listBadges,
    listCategories,
    listContent,
    updateContent
} from './api'
import {useAuth} from './auth'
import {ContentCard} from './components/ContentCard'
import {ContentFilters} from './components/ContentFilters'
import {Shell} from './components/Shell'
import {type Badge, type Category} from './models'

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
            <ContentFilters categories={categories} badgeOptions={badgeOptions} category={category} status={status}
                badges={badges} badgeInput={badgeInput} onCategoryChange={setCategory} onStatusChange={setStatus}
                onBadgesChange={setBadges} onBadgeInputChange={setBadgeInput}/>
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
