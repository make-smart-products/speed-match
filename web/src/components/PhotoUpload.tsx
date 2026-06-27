import { useRef } from 'react'

interface Props {
  onSelect: (file: File | undefined) => void
  preview?: string | null
}

export function PhotoUpload({ onSelect, preview }: Props) {
  const inputRef = useRef<HTMLInputElement>(null)

  return (
    <div className="flex items-center gap-3">
      {preview ? (
        <img src={preview} alt="Превью" className="h-14 w-14 rounded-full object-cover" />
      ) : (
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gray-100 text-gray-400">
          📷
        </div>
      )}
      <div>
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          className="rounded-lg border border-violet-200 px-3 py-2 text-sm font-medium text-violet-700 hover:bg-violet-50"
        >
          Выбрать фото
        </button>
        <input
          ref={inputRef}
          type="file"
          accept="image/jpeg,image/png"
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0]
            onSelect(file)
          }}
        />
        <p className="mt-1 text-xs text-gray-400">JPEG или PNG, до 2 МБ</p>
      </div>
    </div>
  )
}
