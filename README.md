# Profile

Personal portfolio for **Mohd Hujaifa**: a Go API backed by MySQL, Redis, and RabbitMQ, with a React UI.

## Architecture

```
Browser (React / Vite :5173)
        │  JSON over HTTP
        ▼
Go API  Gin  :8080
        ├── GET  /api/v1/portfolio  → Redis cache-aside → MySQL
        └── POST /api/v1/contact    → MySQL insert → worker pool → RabbitMQ
```

- **MySQL** is the source of truth for profile, skills, experience, projects, education, and contact messages.
- **Redis** caches the assembled portfolio payload (`portfolio:v1`) to keep public reads cheap.
- **RabbitMQ** receives `contact.created` events on exchange `portfolio.events`.
- If Redis or RabbitMQ is down, the API still serves: in-memory cache and a no-op publisher.

## Folder structure

```
cmd/api/                 HTTP process
internal/
  cache/                 Redis + in-memory cache
  config/                env-based config
  db/                    schema + migrate
  handler/               Gin routes
  model/                 DTOs
  queue/                 RabbitMQ publisher
  repository/            MySQL access
  seed/                  resume seed data
  service/               business logic
  worker/                goroutine worker pool
frontend/                React SPA
docker-compose.yml       MySQL, Redis, RabbitMQ
```

## Run

```bash
docker compose up -d
cp .env.example .env
go run ./cmd/api
cd frontend && npm install && npm run dev
```

Open http://localhost:5173
