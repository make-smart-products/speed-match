import { getAvatarColor, getInitials } from '../lib/utils'
import type { Participant } from '../api/types'

interface Props {
  participant: Participant
  selected: boolean
  disabled?: boolean
  onToggle: (id: string) => void
}

export function ParticipantCard({ participant, selected, disabled, onToggle }: Props) {
  return (
    <button
      type="button"
      onClick={() => onToggle(participant.id)}
      disabled={disabled}
      className={[
        'flex flex-col items-center gap-2 rounded-2xl border-2 p-3 text-center transition-all',
        'min-h-[120px] w-full touch-manipulation',
        selected
          ? 'border-violet-500 bg-violet-50 shadow-md'
          : 'border-transparent bg-white shadow-sm hover:shadow-md',
        disabled && !selected ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
      ].join(' ')}
    >
      {participant.photo_url ? (
        <img
          src={participant.photo_url}
          alt={participant.pseudonym}
          className="h-16 w-16 rounded-full object-cover"
        />
      ) : (
        <div
          className={`flex h-16 w-16 items-center justify-center rounded-full text-lg font-semibold text-white ${getAvatarColor(participant.pseudonym)}`}
        >
          {getInitials(participant.pseudonym)}
        </div>
      )}
      <span className="text-sm font-medium leading-tight">{participant.pseudonym}</span>
      {selected && (
        <span className="text-xs font-semibold text-violet-600">Выбрано</span>
      )}
    </button>
  )
}
