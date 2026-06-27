import { type FormEvent, useCallback, useEffect, useState } from 'react'
import { Link, useParams, useSearchParams } from 'react-router-dom'
import { api, getAdminToken, saveAdminToken } from '../api/client'
import type { AdminEventResponse } from '../api/types'
import { PhotoUpload } from '../components/PhotoUpload'
import { QrLink } from '../components/QrLink'
import { participantLink } from '../lib/utils'

export function AdminEventPage() {
  const { slug } = useParams<{ slug: string }>()
  const [searchParams] = useSearchParams()
  const [data, setData] = useState<AdminEventResponse | null>(null)
  const [pseudonym, setPseudonym] = useState('')
  const [photo, setPhoto] = useState<File | undefined>()
  const [photoPreview, setPhotoPreview] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [adding, setAdding] = useState(false)
  const [updating, setUpdating] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const keyFromUrl = searchParams.get('key')
  const adminToken = slug
    ? keyFromUrl || getAdminToken(slug)
    : null

  useEffect(() => {
    if (keyFromUrl && slug) saveAdminToken(slug, keyFromUrl)
  }, [keyFromUrl, slug])

  const load = useCallback(() => {
    if (!slug || !adminToken) {
      setLoading(false)
      return
    }
    setLoading(true)
    setError(null)
    api
      .getAdminEvent(slug, adminToken)
      .then(setData)
      .catch((e) => setError(e instanceof Error ? e.message : 'Ошибка загрузки'))
      .finally(() => setLoading(false))
  }, [slug, adminToken])

  useEffect(() => {
    load()
  }, [load])

  async function handleAdd(e: FormEvent) {
    e.preventDefault()
    if (!slug || !adminToken || !pseudonym.trim()) return

    setAdding(true)
    setError(null)
    try {
      await api.addParticipant(slug, adminToken, pseudonym.trim(), photo)
      setPseudonym('')
      setPhoto(undefined)
      setPhotoPreview(null)
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка добавления')
    } finally {
      setAdding(false)
    }
  }

  async function setStatus(status: 'voting' | 'closed' | 'draft') {
    if (!slug || !adminToken) return
    setUpdating(true)
    setError(null)
    try {
      await api.patchEvent(slug, adminToken, { status })
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления')
    } finally {
      setUpdating(false)
    }
  }

  if (!adminToken) {
    return (
      <div className="px-4 py-12 text-center text-red-600">
        Нет доступа. Откройте панель по ссылке с ключом организатора.
      </div>
    )
  }

  if (loading && !data) {
    return <p className="py-12 text-center text-gray-500">Загрузка...</p>
  }

  if (!data) {
    return <p className="py-12 text-center text-red-600">{error || 'Событие не найдено'}</p>
  }

  const { event, participants } = data

  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <Link to="/" className="text-sm text-violet-600 underline">
        ← На главную
      </Link>

      <header className="mt-4 mb-8">
        <h1 className="text-2xl font-bold">{event.title}</h1>
        <p className="text-gray-500">
          Статус:{' '}
          <span className="font-medium text-violet-700">
            {statusLabel(event.status)}
          </span>
          {event.vote_limit != null && (
            <span className="ml-2">· Лимит: {event.vote_limit}</span>
          )}
        </p>
      </header>

      <section className="mb-8 flex flex-wrap gap-3">
        {event.status === 'draft' && (
          <button
            type="button"
            disabled={updating || participants.length < 2}
            onClick={() => setStatus('voting')}
            className="rounded-xl bg-emerald-600 px-4 py-2 font-medium text-white hover:bg-emerald-700 disabled:opacity-50"
          >
            Открыть голосование
          </button>
        )}
        {event.status === 'voting' && (
          <button
            type="button"
            disabled={updating}
            onClick={() => setStatus('closed')}
            className="rounded-xl bg-rose-600 px-4 py-2 font-medium text-white hover:bg-rose-700 disabled:opacity-50"
          >
            Закрыть голосование
          </button>
        )}
        {event.status === 'closed' && (
          <button
            type="button"
            disabled={updating}
            onClick={() => setStatus('voting')}
            className="rounded-xl border border-violet-300 px-4 py-2 font-medium text-violet-700 hover:bg-violet-50 disabled:opacity-50"
          >
            Открыть снова
          </button>
        )}
      </section>

      <section className="mb-8 rounded-2xl bg-white p-6 shadow-sm">
        <h2 className="mb-4 font-semibold">Добавить участника</h2>
        <form onSubmit={handleAdd} className="space-y-4">
          <input
            type="text"
            required
            value={pseudonym}
            onChange={(e) => setPseudonym(e.target.value)}
            placeholder="Псевдоним"
            className="w-full rounded-xl border border-gray-200 px-4 py-3 outline-none focus:border-violet-400"
          />
          <PhotoUpload
            preview={photoPreview}
            onSelect={(file) => {
              setPhoto(file)
              if (file) {
                setPhotoPreview(URL.createObjectURL(file))
              } else {
                setPhotoPreview(null)
              }
            }}
          />
          <button
            type="submit"
            disabled={adding}
            className="rounded-xl bg-violet-600 px-4 py-2 font-medium text-white hover:bg-violet-700 disabled:opacity-50"
          >
            {adding ? 'Добавление...' : 'Добавить'}
          </button>
        </form>
      </section>

      {error && <p className="mb-4 text-sm text-red-600">{error}</p>}

      <section>
        <h2 className="mb-4 font-semibold">
          Участники ({participants.length})
        </h2>
        <div className="space-y-4">
          {participants.map((p) => (
            <QrLink
              key={p.id}
              label={p.pseudonym}
              url={participantLink(event.slug, p.access_token!)}
            />
          ))}
        </div>
      </section>
    </div>
  )
}

function statusLabel(status: string): string {
  switch (status) {
    case 'draft':
      return 'Подготовка'
    case 'voting':
      return 'Голосование открыто'
    case 'closed':
      return 'Завершено'
    default:
      return status
  }
}
