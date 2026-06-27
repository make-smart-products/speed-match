import { Link } from 'react-router-dom'

export function HomePage() {
  return (
    <div className="mx-auto max-w-lg px-4 py-16 text-center">
      <div className="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-3xl bg-violet-600 text-3xl text-white shadow-lg">
        ♥
      </div>
      <h1 className="text-3xl font-bold text-violet-900">Speed Match</h1>
      <p className="mt-3 text-gray-600">
        Анонимные симпатии и взаимные мэтчи для офлайн-мероприятий
      </p>

      <div className="mt-10 space-y-4">
        <Link
          to="/admin/new"
          className="block rounded-2xl bg-violet-600 px-6 py-4 font-semibold text-white shadow-lg hover:bg-violet-700"
        >
          Создать мероприятие
        </Link>
        <p className="text-sm text-gray-400">
          Участники входят по персональной ссылке от организатора
        </p>
      </div>

      <div className="mt-12 rounded-2xl bg-white p-5 text-left text-sm text-gray-600 shadow-sm">
        <p className="font-semibold text-gray-800">Как это работает</p>
        <ol className="mt-2 list-inside list-decimal space-y-1">
          <li>Организатор создаёт событие и добавляет участников</li>
          <li>Участники отмечают симпатии на телефоне</li>
          <li>После закрытия голосования видны только взаимные мэтчи</li>
        </ol>
      </div>
    </div>
  )
}
