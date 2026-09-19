import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { apiConfig, requestSession, signInWithGoogle, sudoLogin } from './api'

// There is no sign-in or sign-out: the upstream SSO proxy decides whether the
// current request carries the owner's identity. Signing in just sends the
// browser to the configured Google/IdP URL.
type AuthState = {
  loading: boolean
  authenticated: boolean
  canEdit: boolean
  username?: string
  error: string | null
  refresh: () => Promise<void>
  canSignIn: boolean
  signIn: () => void
  sudo: (username: string, password: string) => Promise<void>
  signOut: () => void
}

const AuthContext = createContext<AuthState | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authenticated, setAuthenticated] = useState(false)
  const [canEdit, setCanEdit] = useState(false)
  const [username, setUsername] = useState<string | undefined>()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      setError(null)
      const session = await requestSession()
      setAuthenticated(session?.authenticated === true)
      setCanEdit(session?.canEdit === true)
      setUsername(session?.username)
    } catch (err) {
      setAuthenticated(false)
      setCanEdit(false)
      setError(err instanceof Error ? err.message : 'Unable to check your SSO session.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { void refresh() }, [refresh])

  const signIn = useCallback(() => { signInWithGoogle() }, [])
  const sudo = useCallback(async (username: string, password: string) => {
    const result = await sudoLogin(username, password)
    window.sessionStorage.setItem('cms-demo-token', result.token)
    setAuthenticated(true)
    setCanEdit(true)
    setUsername(result.user.username)
  }, [])
  const signOut = useCallback(() => {
    window.sessionStorage.removeItem('cms-demo-token')
    setAuthenticated(false)
    setCanEdit(false)
    setUsername(undefined)
  }, [])
  const canSignIn = apiConfig.signIn !== ''

  const value = useMemo(
    () => ({ loading, authenticated, canEdit, username, error, refresh, canSignIn, signIn, sudo, signOut }),
    [loading, authenticated, canEdit, username, error, refresh, canSignIn, signIn, sudo, signOut],
  )
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used inside AuthProvider')
  return context
}