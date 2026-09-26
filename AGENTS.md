# AGENTS.md

Go (Gin) e-commerce shop. Clean Architecture in the style of the `gapbox` reference repo.

## Layout & dependency direction

```
interfaces (http/worker) → application (usecases/dto/pricing) → infrastructure (repositories/database) → domain (entities/repositories)
```

- `domain/` — entities, repository **interfaces** (`domain/repositories`), `domain_err` (Persian sentinel errors)
- `application/` — usecases (former services), DTOs (`dto/admin`, `dto/web`), `PricingService`
- `infrastructure/` — GORM repository implementations, PostgreSQL, seeders, events, sms/payment, typesense client
- `interfaces/` — Gin handlers, routes, requests, middleware, response renderer, asynq worker tasks
- `bootstrap/` — config load + Dependencies singleton; `cmd/` — cobra (serve, worker, scheduler, migrate, seed, search:reindex)
- `migrations/` — **goose** SQL (PostgreSQL only); `templates/` — Gin HTML (admin/site/layouts/errors)
- `public/` — static assets (`/assets`) and uploads (`/uploads`)

**Never import `infrastructure/...` from `domain/`. Handlers depend on usecase interfaces, not repositories.**

## Commands

| Command | Purpose |
|---|---|
| `go run . serve` / `make run` | HTTP + worker + scheduler in one process |
| `go run . worker` / `make start-worker` | asynq worker only |
| `go run . scheduler` / `make start-schedule` | asynq scheduler only |
| `go run . migrate` / `make migration-up` | goose up |
| `go run . migrate:rollback` / `make migration-down` | goose down (last) |
| `go run . migrate:status` | goose status |
| `go run . migrate:reset` | goose down all (re-run `migrate` after) |
| `go run . make:migration NAME=add_x` | create migration stub |
| `go run . seed` / `make seed` | seed data (needs migrated schema) |
| `go run . search:reindex [--recreate]` | rebuild typesense index from read models |
| `make test` | `go test ./...` (set `TEST_DATABASE_URL` to include integration tests) |
| `go build ./... && go vet ./...` | verify before commit |

## Gotchas

- **Database is PostgreSQL 18 only.** Migrations are goose SQL in `migrations/` with `-- +goose Up` / `-- +goose Down`. Never use AutoMigrate.
- **No MongoDB.** The flattened product lives in `products.read_model` (JSONB), rebuilt by `product.SyncReadModel` after every product/variant mutation. Recommendations live in `product_recommendations`.
- **Variant model** (Laravel-style): stock/prices live on `product_variants`; attribute links on `variant_attribute_values`; after every variant change call `PricingService.RefreshProductAggregates` (aggregates: min/max price, stock, `attributes_json` + GIN, `product_type`). Variant price columns are NULL-able = inherit product price.
- **HTML templates are frozen** — do not edit files under `templates/`. The contracts they rely on: `PRODUCT._id` is the *numeric* product id, inventory map keys `inventory_id/quantity/attributes`, delete link `/admins/product-inventory-attributes/:id/delete`, cart posts `product_id` as decimal string.
- **Typesense is optional at boot** (warn, don't fatal). Rich schema must be recreated after schema changes: `search:reindex --recreate`.
- Uploads write to `public/uploads/...` (config `Upload.*`); image URLs are served from `/uploads`.
- Sessions + middleware are registered in `interfaces/http/routes` via `bootstrap.Initialize()` from `cmd/serve.go`.
- Persian for user-facing messages (`domain/domain_err`, `infrastructure/messages`).

## Tests

- Unit tests run everywhere (`make test`).
- Integration tests **skip** unless `TEST_DATABASE_URL` is set and schema is migrated (`go run . migrate` first). They use transaction rollback — no cleanup needed.

## graphify

A knowledge graph of this repo lives in `graphify-out/` (1,315 nodes · 2,422 edges · 115 communities). Scope: Go source + `templates/` — `public/` and images are excluded.

**Use it before grepping.** For any question about how something works, who calls what, or how data flows across layers, run the query first:

```
graphify query "<question>"          # BFS, broad context (default)
graphify query "<question>" --dfs    # DFS, trace one path
graphify path "AuthMiddleware" "Database"
graphify explain "SyncReadModel"
```

A scoped subgraph is almost always smaller and more relevant than reading `GRAPH_REPORT.md` — save that file for broad architecture context only. Quote `source_location` when citing a specific fact.

**After editing code**, the `post-commit` git hook re-extracts changed files and rebuilds `graph.json` + `GRAPH_REPORT.md` automatically (AST only — template/doc changes need `/graphify --update`). `post-checkout` rebuilds after branch switches, and `graph.json` has a registered merge driver so parallel branch edits merge rather than conflict.

**Known gap:** the graph has no edges linking `templates/*.html` to the Go handlers that render them — import-path nodes are dropped at build time. For handler↔template contracts, read the template and the handler directly.

Rebuild from scratch: `/graphify .` (or `graphify extract --force`). Incremental: `graphify update`.
