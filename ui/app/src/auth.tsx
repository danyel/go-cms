import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { requestSession, signInWithGoogle, signOut, type User } from './api'

type AuthState = {
  user: User | null
  loading: boolean
  error: string | null
  signIn: () => Promise<void>
  signOut: () => Promise<void>
}

const AuthContext = createContext<AuthState | undefined>(undefined)

function getUser(result: { authenticated?: boolean; user?: User; capabilities?: Record<string, boolean>; canEdit?: boolean; token?: string; accessToken?: string; sessionToken?: string }) {
  const token = result.token ?? result.accessToken ?? result.sessionToken
  if (token) sessionStorage.setItem('cms.auth.token', token)
  if (result.authenticated === false) return null
  if (!result.user && !token && result.authenticated !== true) return null
  return { ...(result.user ?? {}), capabilities: { ...result.capabilities, ...(result.canEdit !== undefined ? { canEdit: result.canEdit } : {}), ...result.user?.capabilities } }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      setError(null)
      const session = await requestSession()
      setUser(session ? getUser(session) : null)
    } catch (err) {
      setUser(null)
      setError(err instanceof Error ? err.message : 'Unable to check your session.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { void refresh() }, [refresh])

  const signIn = useCallback(async () => {
    setError(null)
    try {
      const result = await signInWithGoogle()
      setUser(getUser(result))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Sign-in failed.')
    }
  }, [])

  const logOut = useCallback(async () => {
    await signOut().catch(() => undefined)
    sessionStorage.removeItem('cms.auth.token')
    setUser(null)
  }, [])

  const value = useMemo(() => ({ user, loading, error, signIn, signOut: logOut }), [user, loading, error, signIn, logOut])
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used inside AuthProvider')
  return context
}
