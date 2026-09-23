
# 🛍️ Shop

A **scalable, session-based web application** built with **Golang**, using Clean Architecture (gapbox-style layers), **PostgreSQL 18**, **Typesense** realtime search and **asynq** job queues.

📖 Architecture, database and the restructure report live in [`doc/`](doc/):
- [doc/architecture.md](doc/architecture.md) — لایه‌ها و جهت وابستگی
- [doc/database.md](doc/database.md) — schema، مدل واریانت و goose migrations
- [doc/restructure.md](doc/restructure.md) — گزارش کامل تغییرات فازبه‌فاز
- [AGENTS.md](AGENTS.md) — قرارداد اجرا (فرمان‌ها، gotcha ها)

---

## 🚀 Features
- **Clean Architecture**: `domain → application → infrastructure → interfaces` (+ `bootstrap`, `cmd`, `pkg`)
- **PostgreSQL only** — goose SQL migrations (no AutoMigrate, no MySQL, **no MongoDB**)
- Flattened product read model in `products.read_model` (**JSONB**)
- **Laravel-style variant model**: per-combination stock/price on `product_variants`, `attributes_json` (GIN) + `PricingService` aggregates
- **Typesense** rich index with `search:reindex` and graceful degradation when the engine is down
- Background jobs/schedules with **Asynq**
- Unit + DB-gated integration tests

---

## 📋 Prerequisites
- **Go** >= 1.23, **Make**
- **PostgreSQL** (docker: `make dev_server_up` starts postgres + redis + typesense + minio)
- Redis (cache/sessions/asynq), Typesense (search — optional at boot)
- `.env` file:
  ```bash
  cp .env.example .env
  ```
- Proxy/dependencies (if needed):
  ```bash
  go env -w GOPROXY=https://goproxy.io,direct
  go mod download
  ```

---

## 🛠️ Setup and Run

### 1️⃣ Migrate + seed
```bash
make migration-up      # goose up  (go run . migrate)
make seed              # sample data (go run . seed)
```

Other migration commands:
```bash
make migration-status  # goose status
make migration-down    # goose down (last migration)
make migration-reset   # goose down all — run migration-up again after
make make-migration NAME=add_xxx
```

### 2️⃣ Run the application
```bash
make run               # HTTP + worker + scheduler  (go run . serve)
```

Optional standalone processes:
```bash
make start-worker      # go run . worker
make start-schedule    # go run . scheduler
```

Asynq monitoring: `/admins/monitoring` (admin session required).

---

## 🔎 Search (Typesense)

```bash
go run . search:reindex              # re-index all products from read models
go run . search:reindex --recreate   # drop + recreate the collection schema first
```
Search UI: `/tsearch` and `/tsearch/show`. The engine is optional at boot — if it is down the app still starts and search returns a graceful error.

---

## 🧪 Tests

```bash
make test        # unit + handler tests (integration tests skip without DB)
make test-cover

# with a migrated database, integration tests run too:
TEST_DATABASE_URL=postgres://user:pass@127.0.0.1:5432/shop?sslmode=disable make test
```

---

## 🧑‍💻 Contributing
1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. `go build ./... && go vet ./... && make test`
4. Commit and open a pull request.

---

## 📜 License
MIT — see [LICENSE](LICENSE).

---

🎉 Happy coding!
