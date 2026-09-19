import {useEffect, useState, type ReactNode} from 'react'
import {Link} from 'react-router-dom'
import {listDemoUsers} from '../api'
import {useAuth} from '../auth'
import {type DemoUser} from '../models'
import {SudoTerminal} from './SudoTerminal'

export function Shell({children}: {children: ReactNode}) {
  const {authenticated, loading, canSignIn, signIn, username, sudo, signOut} = useAuth()
  const [sudoOpen, setSudoOpen] = useState(false)
  const [helpOpen, setHelpOpen] = useState(false)
  const [users, setUsers] = useState<DemoUser[]>([])
  const [dark, setDark] = useState(() => window.localStorage.getItem('urpi-theme') !== 'light')

  useEffect(() => { void listDemoUsers().then(setUsers) }, [])
  useEffect(() => {
    document.documentElement.dataset.theme = dark ? 'dark' : 'light'
    window.localStorage.setItem('urpi-theme', dark ? 'dark' : 'light')
  }, [dark])

  return <main className="shell">
    <header className="header">
      <Link className="brand" to="/content"><span className="brand-prompt">~/</span> Urpi's backlog
        <span className="brand-cursor">_</span><small>linux · go · java · ideas</small></Link>
      <div className="header-actions">
        {!loading && <span className="mode-badge">{authenticated ? `${username ?? 'Owner'} · sudo` : 'Anonymous · read-only'}</span>}
        <button className="theme-button" type="button" aria-label="Demo users"
          onClick={() => setHelpOpen(value => !value)}>?</button>
        {!loading && !authenticated && users.length > 0 &&
          <button className="theme-button" type="button" onClick={() => setSudoOpen(true)}>sudo</button>}
        {!loading && authenticated && users.length > 0 &&
          <button className="theme-button" type="button" onClick={signOut}>logout</button>}
        {!loading && !authenticated && canSignIn &&
          <button className="google-button session-button" type="button" onClick={signIn}>
            <span aria-hidden="true">G</span> Sign in with Google</button>}
        <button className="theme-button" type="button" aria-label={`Switch to ${dark ? 'light' : 'dark'} theme`}
          onClick={() => setDark(value => !value)}>{dark ? '☀ Light' : '☾ Dark'}</button>
      </div>
    </header>
    {helpOpen && users.length > 0 && <aside className="card terminal-panel">
      <strong>sudo users</strong><p className="muted">Demo only: each password is the same as the username.</p>
      {users.map(item => <div key={item.username}><code>sudo {item.username}</code> · {item.name}</div>)}
    </aside>}
    {sudoOpen && <SudoTerminal users={users} onLogin={async (user, password) => {
      await sudo(user, password)
      setSudoOpen(false)
    }} onClose={() => setSudoOpen(false)}/>}
    {children}
  </main>
}
