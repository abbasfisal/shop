run:
	@go run cmd/http/main.go

migrate:
	@go run cmd/http/main.go migrate

seed:
	@go run cmd/http/main.go seed

#run dev server docker compose with specific env file up
dev_server_up:
	@docker-compose -f docker-compose.dev.yml --env-file .env.development up -d

#remove dev docker-compose running server
dev_server_down:
	@docker-compose -f docker-compose.dev.yml down

#run production server docker compose with specific env file up
production_server_up:
	@docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

#remove production docker-compose running server
production_server_down:
	@docker-compose -f docker-compose.prod.yml down