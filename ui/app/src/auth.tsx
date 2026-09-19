import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { apiConfig, requestSession, signInWithGoogle } from './api'

// There is no sign-in or sign-out: the upstream SSO proxy decides whether the
// current request carries the owner's identity. Signing in just sends the
// browser to the configured Google/IdP URL.
type AuthState = {
  loading: boolean
  authenticated: boolean
  canEdit: boolean
  error: string | null
  refresh: () => Promise<void>
  canSignIn: boolean
  signIn: () => void
}

const AuthContext = createContext<AuthState | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authenticated, setAuthenticated] = useState(false)
  const [canEdit, setCanEdit] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      setError(null)
      const session = await requestSession()
      setAuthenticated(session?.authenticated === true)
      setCanEdit(session?.canEdit === true)
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
  const canSignIn = apiConfig.signIn !== ''

  const value = useMemo(
    () => ({ loading, authenticated, canEdit, error, refresh, canSignIn, signIn }),
    [loading, authenticated, canEdit, error, refresh, canSignIn, signIn],
  )
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used inside AuthProvider')
  return context
}