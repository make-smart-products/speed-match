import type {
  AdminEventResponse,
  CreateEventResponse,
  Event,
  EventStatusResponse,
  Participant,
} from './types'

const API = '/api/v1'

class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  headers: Record<string, string> = {},
): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    ...options,
    headers: {
      ...headers,
      ...(options.headers as Record<string, string>),
    },
  })

  if (!res.ok) {
    let msg = res.statusText
    try {
      const body = await res.json()
      if (body.error) msg = body.error
    } catch {
      /* ignore */
    }
    throw new ApiError(msg, res.status)
  }

  if (res.status === 204) return undefined as T
  return res.json()
}

export function saveAccessToken(slug: string, token: string) {
  sessionStorage.setItem(`access:${slug}`, token)
}

export function getAccessToken(slug: string): string | null {
  return sessionStorage.getItem(`access:${slug}`)
}

export function saveAdminToken(slug: string, token: string) {
  sessionStorage.setItem(`admin:${slug}`, token)
}

export function getAdminToken(slug: string): string | null {
  return sessionStorage.getItem(`admin:${slug}`)
}

export const api = {
  getStatus(slug: string, accessToken: string) {
    return request<EventStatusResponse>(`/events/${slug}/status`, {}, {
      'X-Access-Token': accessToken,
    })
  },

  getParticipants(slug: string, accessToken: string) {
    return request<Participant[]>(`/events/${slug}/participants`, {}, {
      'X-Access-Token': accessToken,
    })
  },

  saveVotes(accessToken: string, targetIds: string[]) {
    return request<{ status: string }>(
      '/votes',
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ target_ids: targetIds }),
      },
      { 'X-Access-Token': accessToken },
    )
  },

  getMatches(slug: string, accessToken: string) {
    return request<Participant[]>(`/events/${slug}/matches`, {}, {
      'X-Access-Token': accessToken,
    })
  },

  createEvent(title: string, voteLimit: number | null) {
    return request<CreateEventResponse>('/admin/events', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title,
        vote_limit: voteLimit,
      }),
    })
  },

  getAdminEvent(slug: string, adminToken: string) {
    return request<AdminEventResponse>(`/admin/events/${slug}`, {}, {
      'X-Admin-Token': adminToken,
    })
  },

  addParticipant(slug: string, adminToken: string, pseudonym: string, photo?: File) {
    const form = new FormData()
    form.append('pseudonym', pseudonym)
    if (photo) form.append('photo', photo)

    return request<Participant>(
      `/admin/events/${slug}/participants`,
      { method: 'POST', body: form },
      { 'X-Admin-Token': adminToken },
    )
  },

  patchEvent(slug: string, adminToken: string, data: { status?: Event['status']; vote_limit?: number | null }) {
    return request<Event>(
      `/admin/events/${slug}`,
      {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      },
      { 'X-Admin-Token': adminToken },
    )
  },
}

export { ApiError }
