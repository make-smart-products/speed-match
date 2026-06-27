export type EventStatus = 'draft' | 'voting' | 'closed'

export interface Participant {
  id: string
  event_id: string
  pseudonym: string
  photo_url?: string | null
  access_token?: string
}

export interface EventStatusResponse {
  title: string
  slug: string
  status: EventStatus
  vote_limit: number | null
  selected_ids: string[]
  pseudonym: string
}

export interface Event {
  id: string
  title: string
  slug: string
  vote_limit: number | null
  status: EventStatus
}

export interface CreateEventResponse {
  event: Event
  admin_token: string
  admin_url: string
}

export interface AdminEventResponse {
  event: Event
  participants: Participant[]
}
