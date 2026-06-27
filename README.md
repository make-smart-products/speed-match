# Speed Match

Инструмент для офлайн-мероприятий быстрых знакомств: участники анонимно отмечают симпатии, после закрытия голосования видят только **взаимные мэтчи**. Односторонние симпатии никому не раскрываются.

**Репозиторий:** https://github.com/make-smart-products/speed-match

## Стек

- **Backend:** Go (chi, sqlx, SQLite, goose)
- **Frontend:** React + Vite + TypeScript + Tailwind
- **Деплой:** Docker (один контейнер: API + фронтенд)

## Быстрый старт (Docker — рекомендуется)

```bash
docker compose up --build
```

Откройте **http://localhost:8080** — полноценный веб-сервис (создание мероприятий, голосование, мэтчи).

Данные сохраняются в Docker volume `speedmatch-data`.

## Локальная разработка

### Backend

```bash
cd backend
go run ./cmd/server
```

API: `http://localhost:8080/api/v1`

### Frontend (с hot reload)

```bash
cd web
npm install
npm run dev
```

UI: `http://localhost:5173` (прокси `/api` и `/uploads` на backend)

### Демо-данные

```bash
cd backend
go run ./cmd/seed
```

## Сценарий организатора

1. Откройте `/admin/new`
2. Создайте мероприятие (название, опциональный лимит симпатий)
3. Добавьте участников (псевдоним + опциональное фото)
4. Раздайте персональные QR-ссылки
5. Нажмите **«Открыть голосование»**
6. После мероприятия — **«Закрыть голосование»**

## Сценарий участника

1. Открыть персональную ссылку `/e/{slug}?t={token}`
2. На экране **«Выбор симпатий»** отметить понравившихся людей
3. Нажать **«Сохранить выбор»**
4. После закрытия — экран **«Твои мэтчи»**

## Деплой в облако (Render)

1. Форкните или подключите репозиторий на [Render](https://render.com)
2. **New → Blueprint** и укажите `render.yaml` из репозитория  
   или **New → Web Service → Docker** с этим репозиторием
3. Подключите **Persistent Disk** на `/data` (1 GB) — для БД и фото
4. После деплоя откройте выданный URL (например `https://speed-match.onrender.com`)

Health check: `GET /health`

## API

| Метод | Endpoint | Заголовок | Описание |
|-------|----------|-----------|----------|
| GET | `/health` | — | Проверка работоспособности |
| GET | `/api/v1/events/{slug}/status` | `X-Access-Token` | Статус + текущий выбор |
| GET | `/api/v1/events/{slug}/participants` | `X-Access-Token` | Список других участников |
| POST | `/api/v1/votes` | `X-Access-Token` | Сохранить симпатии |
| GET | `/api/v1/events/{slug}/matches` | `X-Access-Token` | Взаимные мэтчи |
| POST | `/api/v1/admin/events` | — | Создать событие |
| GET | `/api/v1/admin/events/{slug}` | `X-Admin-Token` | Панель организатора |
| POST | `/api/v1/admin/events/{slug}/participants` | `X-Admin-Token` | Добавить участника |
| PATCH | `/api/v1/admin/events/{slug}` | `X-Admin-Token` | Статус / лимит |

## Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `PORT` | `8080` | Порт сервера |
| `DB_PATH` | `data/speedmatch.db` | Путь к SQLite |
| `UPLOAD_DIR` | `uploads` | Папка для фото |
| `STATIC_DIR` | `../web/dist` | Статика фронтенда |
| `CORS_ORIGIN` | `*` в prod | Origin для CORS |

## Production (без Docker)

```bash
cd web && npm ci && npm run build
cd ../backend && go build -o bin/server ./cmd/server
STATIC_DIR=../web/dist DB_PATH=data/speedmatch.db ./bin/server
```

## Безопасность

- Токены участников и организатора — криптостойкие случайные строки
- Нельзя голосовать за себя
- Голоса только при `status=voting`
- Мэтчи только при `status=closed`
- Фото: JPEG/PNG, max 2 МБ
