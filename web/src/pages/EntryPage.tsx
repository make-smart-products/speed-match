import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams, useSearchParams } from 'react-router-dom'
import { api, getAccessToken, saveAccessToken } from '../api/client'
import type { EventStatusResponse } from '../api/types'

export function EntryPage() {
  const { slug } = useParams<{ slug: string }>()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!slug) return

    const tokenFromUrl = searchParams.get('t')
    const token = tokenFromUrl || getAccessToken(slug)

    if (!token) {
      setError('Ссылка недействительна. Попросите организатора персональную ссылку.')
      return
    }

    if (tokenFromUrl) {
      saveAccessToken(slug, tokenFromUrl)
    }

    api
      .getStatus(slug, token)
      .then((status: EventStatusResponse) => {
        if (status.status === 'closed') {
          navigate(`/e/${slug}/results`, { replace: true })
        } else {
          navigate(`/e/${slug}/vote`, { replace: true })
        }
      })
      .catch(() => setError('Не удалось войти. Проверьте ссылку.'))
  }, [slug, searchParams, navigate])

  if (error) {
    return (
      <div className="mx-auto max-w-md px-4 py-16 text-center">
        <p className="text-red-600">{error}</p>
        <Link to="/" className="mt-4 inline-block text-violet-600 underline">
          На главную
        </Link>
      </div>
    )
  }

  return (
    <div className="flex min-h-[50vh] items-center justify-center">
      <p className="text-gray-500">Загрузка...</p>
    </div>
  )
}
