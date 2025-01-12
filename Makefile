include .env

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


date:
	@date +%Y%m%d%H%M%S


migration-up:
	@sql-migrate up -env=production -config=internal/database/mysql/dbconfig.yml


migration-down:
	@sql-migrate down -env=production -config=internal/database/mysql/dbconfig.yml -limit=1

migration-status:
	@sql-migrate status -env=production -config=internal/database/mysql/dbconfig.yml


#doc: https://github.com/golang-migrate/migrate
#generate dbconfig.yml
# this dbconfig.yml file is necessary when you want use sql-migrator pkg to migration-up and down
generate-sql-migrator-dbconfig:
	@echo "production:\
           \n  dialect: mysql\
           \n  datasource: ${MYSQL_USER}:${MYSQL_PASSWORD}@(${MYSQL_HOSTNAME}:${MYSQL_PORT})/${MYSQL_DB}?parseTime=true\
           \n  dir: internal/database/mysql/migrations #migration director\
           \n  table: migrations" > internal/database/mysql/dbconfig.yml


# start schedule system using asynq pkg
start-schedule:
	@go run ./cmd/job/scheduler/main.go


# start worker which is responsible to execute tasks
start-worker:
	@go run ./cmd/job/worker/main.go


# start minio server
minio_run :
	@docker run -p 9000:9000 -p 9001:9001 --name my_golang_minio \
  			-v /data:/data \
  			-e "MINIO_ROOT_USER=vivify" \
  			-e "MINIO_ROOT_PASSWORD=vivify" \
  			minio/minio:RELEASE.2024-05-01T01-11-10Z server /data --console-address ":9001"
