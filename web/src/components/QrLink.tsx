import { QRCodeSVG } from 'qrcode.react'
import { useState } from 'react'

interface Props {
  url: string
  label: string
}

export function QrLink({ url, label }: Props) {
  const [copied, setCopied] = useState(false)

  async function copy() {
    await navigator.clipboard.writeText(url)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="flex flex-col items-center gap-3 rounded-xl bg-white p-4 shadow-sm sm:flex-row sm:items-start">
      <QRCodeSVG value={url} size={96} className="shrink-0 rounded-lg" />
      <div className="min-w-0 flex-1 text-center sm:text-left">
        <p className="font-medium">{label}</p>
        <p className="mt-1 break-all text-xs text-gray-500">{url}</p>
        <button
          type="button"
          onClick={copy}
          className="mt-2 rounded-lg bg-violet-100 px-3 py-1.5 text-sm font-medium text-violet-700 hover:bg-violet-200"
        >
          {copied ? 'Скопировано!' : 'Копировать ссылку'}
        </button>
      </div>
    </div>
  )
}
