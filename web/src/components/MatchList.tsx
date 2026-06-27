import { getAvatarColor, getInitials } from '../lib/utils'
import type { Participant } from '../api/types'

interface Props {
  matches: Participant[]
}

export function MatchList({ matches }: Props) {
  if (matches.length === 0) {
    return (
      <div className="rounded-2xl bg-white p-8 text-center shadow-sm">
        <p className="text-lg text-gray-600">В этот раз мэтчей нет</p>
        <p className="mt-2 text-sm text-gray-400">
          Возможно, в следующий раз повезёт больше
        </p>
      </div>
    )
  }

  return (
    <ul className="grid grid-cols-1 gap-3 sm:grid-cols-2">
      {matches.map((m) => (
        <li
          key={m.id}
          className="flex items-center gap-4 rounded-2xl bg-white p-4 shadow-sm"
        >
          {m.photo_url ? (
            <img
              src={m.photo_url}
              alt={m.pseudonym}
              className="h-14 w-14 rounded-full object-cover"
            />
          ) : (
            <div
              className={`flex h-14 w-14 items-center justify-center rounded-full text-base font-semibold text-white ${getAvatarColor(m.pseudonym)}`}
            >
              {getInitials(m.pseudonym)}
            </div>
          )}
          <div>
            <p className="font-semibold">{m.pseudonym}</p>
            <p className="text-sm text-violet-600">Взаимная симпатия</p>
          </div>
        </li>
      ))}
    </ul>
  )
}
