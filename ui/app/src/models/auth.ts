export type User = {
  id?: string
  name?: string
  email?: string
  picture?: string
  capabilities?: Record<string, boolean>
}

export type SessionResponse = {
  authenticated?: boolean
  canEdit?: boolean
  username?: string
}

export type DemoUser = { username: string; name: string }
