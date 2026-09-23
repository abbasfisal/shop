include .env

run:
	@go run . serve

migrate:
	@go run . migrate

seed:
	@go run . seed

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


date:
	@date +%Y%m%d%H%M%S


migration-up:
	@go run . migrate

migration-down:
	@go run . migrate:rollback

migration-status:
	@go run . migrate:status

migration-reset:
	@go run . migrate:reset

make-migration:
	@go run . make:migration NAME=$(NAME)


# start schedule system using asynq pkg
start-schedule:
	@go run . scheduler


# start worker which is responsible to execute tasks
start-worker:
	@go run . worker


# start minio server
minio_run :
	@docker run -p 9000:9000 -p 9001:9001 --name my_golang_minio \
  			-v /data:/data \
  			-e "MINIO_ROOT_USER=vivify" \
  			-e "MINIO_ROOT_PASSWORD=vivify" \
  			minio/minio:RELEASE.2024-05-01T01-11-10Z server /data --console-address ":9001"
