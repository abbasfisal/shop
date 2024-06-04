run:
	@go run cmd/http/main.go

migrate:
	@go run cmd/http/main.go migrate

seed:
	@go run cmd/http/main.go seed