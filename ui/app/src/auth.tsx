import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { requestSession } from './api'

// There is no sign-in or sign-out: the upstream SSO proxy decides whether the
// current request carries the owner's identity.
type AuthState = {
  loading: boolean
  authenticated: boolean
  canEdit: boolean
  error: string | null
  refresh: () => Promise<void>
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

  const value = useMemo(() => ({ loading, authenticated, canEdit, error, refresh }), [loading, authenticated, canEdit, error, refresh])
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used inside AuthProvider')
  return context
}