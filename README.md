# Speed Match

Инструмент для офлайн-мероприятий быстрых знакомств: участники анонимно отмечают симпатии, после закрытия голосования видят только **взаимные мэтчи**. Односторонние симпатии никому не раскрываются.

**Репозиторий:** https://github.com/make-smart-products/speed-match  
**Демо (Render):** https://speed-match-iw5t.onrender.com

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

## Локальный запуск без Docker

### Вариант A: разработка (hot reload фронтенда)

Два терминала:

```bash
# Терминал 1 — API
cd backend
go run ./cmd/server

# Терминал 2 — фронтенд
cd web
npm install
npm run dev
```

- API: `http://localhost:8080/api/v1`
- UI: `http://localhost:5173` (прокси `/api` и `/uploads` на backend)

### Вариант B: как на проде (один порт)

```bash
cd web && npm ci && npm run build
cd ../backend && go build -o bin/server ./cmd/server
```

**Linux / macOS:**
```bash
STATIC_DIR=../web/dist ./bin/server
```

**Windows (PowerShell):**
```powershell
$env:STATIC_DIR="../web/dist"
.\bin\server.exe
```

Откройте **http://localhost:8080**.

### Демо-данные

```bash
cd backend
go run ./cmd/seed
```

Команда выведет ссылки организатора и участников.

## Сценарий организатора

1. Откройте `/admin/new`
2. Создайте мероприятие (название, опциональный лимит симпатий)
3. Добавьте участников (псевдоним + опциональное фото)
4. Раздайте персональные QR-ссылки
5. Нажмите **«Открыть голосование»** (нужно минимум 2 участника)
6. После мероприятия — **«Закрыть голосование»**

Ссылка на панель организатора: `/admin/{slug}?key={admin_token}`

## Сценарий участника

1. Открыть персональную ссылку `/e/{slug}?t={token}`
2. На экране **«Выбор симпатий»** отметить понравившихся людей
3. Нажать **«Сохранить выбор»**
4. После закрытия — экран **«Твои мэтчи»**

## Деплой в облако (Render)

1. Подключите репозиторий на [Render](https://render.com)
2. **New → Blueprint** и укажите `render.yaml`  
   или **New → Web Service → Docker** с этим репозиторием
3. Persistent Disk на `/data` (1 GB) — для БД и фото (уже в `render.yaml`)
4. Автодеплой при push в `main`

**Передеплой вручную:** Render Dashboard → сервис → **Manual Deploy → Deploy latest commit**

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
| `STATIC_DIR` | `../web/dist` | Статика фронтенда (prod) |
| `CORS_ORIGIN` | `*` в prod | Origin для CORS |

## Безопасность

- Токены участников и организатора — криптостойкие случайные строки
- Нельзя голосовать за себя
- Голоса только при `status=voting`
- Мэтчи только при `status=closed`
- Фото: JPEG/PNG, max 2 МБ

## Make-команды

```bash
make docker    # docker compose up --build
make backend   # go run ./cmd/server
make frontend  # npm run dev
make build     # собрать web + backend
make seed      # демо-данные
```
