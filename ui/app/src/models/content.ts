export type Content = {
  id: string
  title: string
  summary: string
  body: string
  status: string
  category?: string
  badges?: string[]
  createdAt?: string
  updatedAt?: string
  author?: string
  publishedAt?: string
}

export type ContentUpdate = Pick<Content, 'title' | 'summary' | 'body' | 'status' | 'badges'>

export type ContentListResponse = {
  items?: Content[]
  content?: Content[]
  data?: Content[]
  nextCursor?: string | null
  cursor?: string | null
  hasMore?: boolean
  offset?: number
  limit?: number
}

export type Category = { id: string; slug: string; name: string }
export type Badge = { id: string; name: string }
