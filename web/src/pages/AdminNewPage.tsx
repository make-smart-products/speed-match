import { type FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api, saveAdminToken } from '../api/client'

export function AdminNewPage() {
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [voteLimit, setVoteLimit] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const limit = voteLimit.trim() === '' ? null : parseInt(voteLimit, 10)
      if (limit !== null && (isNaN(limit) || limit < 1)) {
        setError('Лимит должен быть положительным числом')
        setLoading(false)
        return
      }

      const resp = await api.createEvent(title.trim(), limit)
      saveAdminToken(resp.event.slug, resp.admin_token)
      navigate(`/admin/${resp.event.slug}?key=${resp.admin_token}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-md px-4 py-12">
      <h1 className="mb-2 text-2xl font-bold">Новое мероприятие</h1>
      <p className="mb-8 text-gray-500">Создайте событие и добавьте участников</p>

      <form onSubmit={handleSubmit} className="space-y-4 rounded-2xl bg-white p-6 shadow-sm">
        <div>
          <label className="mb-1 block text-sm font-medium">Название</label>
          <input
            type="text"
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Быстрые знакомства — весна 2026"
            className="w-full rounded-xl border border-gray-200 px-4 py-3 outline-none focus:border-violet-400"
          />
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium">
            Лимит симпатий <span className="text-gray-400">(пусто = без ограничений)</span>
          </label>
          <input
            type="number"
            min={1}
            value={voteLimit}
            onChange={(e) => setVoteLimit(e.target.value)}
            placeholder="Не ограничено"
            className="w-full rounded-xl border border-gray-200 px-4 py-3 outline-none focus:border-violet-400"
          />
        </div>

        {error && <p className="text-sm text-red-600">{error}</p>}

        <button
          type="submit"
          disabled={loading}
          className="w-full min-h-[44px] rounded-xl bg-violet-600 py-3 font-semibold text-white hover:bg-violet-700 disabled:opacity-50"
        >
          {loading ? 'Создание...' : 'Создать'}
        </button>
      </form>
    </div>
  )
}
