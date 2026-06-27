import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { api, getAccessToken } from '../api/client'
import type { EventStatusResponse, Participant } from '../api/types'
import { ParticipantGrid } from '../components/ParticipantGrid'

export function VotePage() {
  const { slug } = useParams<{ slug: string }>()
  const navigate = useNavigate()
  const [status, setStatus] = useState<EventStatusResponse | null>(null)
  const [participants, setParticipants] = useState<Participant[]>([])
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const token = slug ? getAccessToken(slug) : null

  useEffect(() => {
    if (!slug || !token) {
      navigate(`/e/${slug}?t=`)
      return
    }

    Promise.all([api.getStatus(slug, token), api.getParticipants(slug, token)])
      .then(([st, list]) => {
        if (st.status === 'closed') {
          navigate(`/e/${slug}/results`, { replace: true })
          return
        }
        if (st.status !== 'voting') {
          setError('Голосование ещё не открыто. Подождите организатора.')
          return
        }
        setStatus(st)
        setParticipants(list)
        setSelected(new Set(st.selected_ids))
      })
      .catch((e) => setError(e.message || 'Ошибка загрузки'))
      .finally(() => setLoading(false))
  }, [slug, token, navigate])

  const toggle = useCallback(
    (id: string) => {
      setSelected((prev) => {
        const next = new Set(prev)
        if (next.has(id)) {
          next.delete(id)
        } else {
          if (status?.vote_limit !== null && status?.vote_limit !== undefined && next.size >= status.vote_limit) {
            return prev
          }
          next.add(id)
        }
        return next
      })
      setSaved(false)
    },
    [status],
  )

  async function handleSave() {
    if (!slug || !token) return
    setSaving(true)
    setError(null)
    try {
      await api.saveVotes(token, Array.from(selected))
      setSaved(true)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Не удалось сохранить')
    } finally {
      setSaving(false)
    }
  }

  const overLimit =
    status?.vote_limit != null && selected.size > status.vote_limit

  if (loading) {
    return <Loading />
  }

  if (error && !status) {
    return (
      <div className="mx-auto max-w-lg px-4 py-12 text-center">
        <p className="text-gray-600">{error}</p>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-4xl px-4 pb-28 pt-6">
      <header className="mb-6">
        <p className="text-sm text-violet-600">Привет, {status?.pseudonym}!</p>
        <h1 className="text-2xl font-bold">{status?.title}</h1>
        <p className="mt-1 text-gray-500">Отметьте людей, которые вам понравились</p>
      </header>

      <ParticipantGrid
        participants={participants}
        selectedIds={selected}
        voteLimit={status?.vote_limit ?? null}
        onToggle={toggle}
      />

      {error && <p className="mt-4 text-center text-sm text-red-600">{error}</p>}
      {saved && <p className="mt-4 text-center text-sm text-emerald-600">Выбор сохранён</p>}

      <div className="fixed bottom-0 left-0 right-0 border-t border-violet-100 bg-white/95 p-4 backdrop-blur">
        <div className="mx-auto flex max-w-4xl items-center justify-between gap-4">
          <div className="text-sm">
            <span className="font-semibold">Выбрано: {selected.size}</span>
            {status?.vote_limit != null && (
              <span className="text-gray-500"> / {status.vote_limit}</span>
            )}
            {overLimit && (
              <p className="text-xs text-red-500">Превышен лимит</p>
            )}
          </div>
          <button
            type="button"
            onClick={handleSave}
            disabled={saving || overLimit}
            className="min-h-[44px] rounded-xl bg-violet-600 px-6 py-3 font-semibold text-white shadow-lg transition hover:bg-violet-700 disabled:opacity-50"
          >
            {saving ? 'Сохранение...' : 'Сохранить выбор'}
          </button>
        </div>
      </div>

      {status?.status === 'voting' && (
        <p className="mt-6 text-center text-sm text-gray-400">
          Результаты появятся после завершения голосования.{' '}
          <Link to={`/e/${slug}/results`} className="text-violet-600 underline">
            Проверить
          </Link>
        </p>
      )}
    </div>
  )
}

function Loading() {
  return (
    <div className="flex min-h-[50vh] items-center justify-center">
      <p className="text-gray-500">Загрузка...</p>
    </div>
  )
}
