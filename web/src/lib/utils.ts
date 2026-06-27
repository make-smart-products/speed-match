const COLORS = [
  'bg-rose-400',
  'bg-violet-400',
  'bg-sky-400',
  'bg-emerald-400',
  'bg-amber-400',
  'bg-fuchsia-400',
  'bg-teal-400',
  'bg-orange-400',
]

export function getInitials(name: string): string {
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.slice(0, 2).toUpperCase()
}

export function getAvatarColor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return COLORS[Math.abs(hash) % COLORS.length]
}

export function participantLink(slug: string, token: string): string {
  const origin = window.location.origin
  return `${origin}/e/${slug}?t=${token}`
}

export function adminLink(slug: string, token: string): string {
  const origin = window.location.origin
  return `${origin}/admin/${slug}?key=${token}`
}
