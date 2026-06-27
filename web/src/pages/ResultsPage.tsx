import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { api, getAccessToken } from '../api/client'
import type { EventStatusResponse, Participant } from '../api/types'
import { MatchList } from '../components/MatchList'

export function ResultsPage() {
  const { slug } = useParams<{ slug: string }>()
  const navigate = useNavigate()
  const [status, setStatus] = useState<EventStatusResponse | null>(null)
  const [matches, setMatches] = useState<Participant[] | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const token = slug ? getAccessToken(slug) : null

  useEffect(() => {
    if (!slug || !token) {
      navigate(`/e/${slug}`)
      return
    }

    api
      .getStatus(slug, token)
      .then(async (st) => {
        setStatus(st)
        if (st.status !== 'closed') {
          setMatches(null)
          return
        }
        const list = await api.getMatches(slug, token)
        setMatches(list)
      })
      .catch((e) => setError(e.message || 'Ошибка загрузки'))
      .finally(() => setLoading(false))
  }, [slug, token, navigate])

  if (loading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <p className="text-gray-500">Загрузка...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="mx-auto max-w-lg px-4 py-12 text-center">
        <p className="text-red-600">{error}</p>
      </div>
    )
  }

  if (status?.status !== 'closed') {
    return (
      <div className="mx-auto max-w-lg px-4 py-16 text-center">
        <div className="rounded-2xl bg-white p-8 shadow-sm">
          <p className="text-4xl">⏳</p>
          <h1 className="mt-4 text-xl font-bold">Результаты скоро</h1>
          <p className="mt-2 text-gray-500">
            Результаты появятся после завершения голосования
          </p>
          <Link
            to={`/e/${slug}/vote`}
            className="mt-6 inline-block text-violet-600 underline"
          >
            Вернуться к выбору
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <header className="mb-6 text-center">
        <h1 className="text-2xl font-bold">Твои мэтчи</h1>
        <p className="mt-1 text-gray-500">{status.title}</p>
      </header>

      <MatchList matches={matches ?? []} />
    </div>
  )
}
