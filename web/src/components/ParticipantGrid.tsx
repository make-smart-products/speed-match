import { useMemo, useState } from 'react'
import type { Participant } from '../api/types'
import { ParticipantCard } from './ParticipantCard'

interface Props {
  participants: Participant[]
  selectedIds: Set<string>
  voteLimit: number | null
  onToggle: (id: string) => void
}

export function ParticipantGrid({ participants, selectedIds, voteLimit, onToggle }: Props) {
  const [query, setQuery] = useState('')

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return participants
    return participants.filter((p) => p.pseudonym.toLowerCase().includes(q))
  }, [participants, query])

  const atLimit = voteLimit !== null && selectedIds.size >= voteLimit

  return (
    <div className="space-y-4">
      <div className="relative">
        <input
          type="search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Поиск по псевдониму..."
          className="w-full rounded-xl border border-violet-100 bg-white px-4 py-3 pl-10 text-base shadow-sm outline-none focus:border-violet-400"
        />
        <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-violet-300">
          ⌕
        </span>
      </div>

      {filtered.length === 0 ? (
        <p className="py-8 text-center text-gray-500">Никого не найдено</p>
      ) : (
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
          {filtered.map((p) => (
            <ParticipantCard
              key={p.id}
              participant={p}
              selected={selectedIds.has(p.id)}
              disabled={atLimit}
              onToggle={onToggle}
            />
          ))}
        </div>
      )}
    </div>
  )
}
