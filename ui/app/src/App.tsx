import { Navigate, Route, Routes, Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from './auth'

function Shell({ children, admin = false }: { children: React.ReactNode; admin?: boolean }) {
  const { user, signOut } = useAuth()
  return <main className="shell">
    <header className="header"><Link className="brand" to="/">CMS</Link>{user && !admin && <button className="text-button" onClick={() => void signOut()}>Sign out</button>}</header>
    {children}
  </main>
}

function Login({ admin = false }: { admin?: boolean }) {
  const { user, loading, error, signIn } = useAuth()
  const navigate = useNavigate()
  if (!admin && user) return <Navigate to="/test" replace />
  if (loading) return <Shell admin={admin}><div className="card"><p>Checking your session…</p></div></Shell>
  return <Shell admin={admin}><section className="card auth-card">
    <span className="eyebrow">{admin ? 'Administration' : 'Welcome'}</span>
    <h1>{admin ? 'Admin sign in' : 'Sign in to CMS'}</h1>
    <p className="muted">{admin ? 'Use your administrator account to continue.' : 'Sign in securely with your external identity provider.'}</p>
    {error && <p className="error" role="alert">{error}</p>}
    {!admin && <button className="google-button" onClick={() => void signIn()}><span aria-hidden="true">G</span> Continue with Google</button>}
    {admin && <p className="muted small">Administrator authentication is handled by the backend.</p>}
    {admin && <button className="secondary-button" onClick={() => navigate('/')}>Return to home</button>}
  </section></Shell>
}

function ProtectedTest() {
  const { user, loading } = useAuth()
  const location = useLocation()
  if (loading) return <Shell><div className="card"><p>Loading…</p></div></Shell>
  if (!user) return <Navigate to="/" state={{ from: location }} replace />
  return <Shell><section className="card"><span className="eyebrow">Authenticated</span><h1>Test page</h1><p>You are signed in and can access protected content.</p><div className="profile">{user.email ?? user.name ? `Signed in as ${user.email ?? user.name}` : 'Your session is active.'}</div></section></Shell>
}

export function App() {
  return <Routes>
    <Route path="/" element={<Login />} />
    <Route path="/admin/login" element={<Login admin />} />
    <Route path="/test" element={<ProtectedTest />} />
    <Route path="*" element={<Navigate to="/" replace />} />
  </Routes>
}
