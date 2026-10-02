include .env
export

run:
	@go run cmd/financer/main.go

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

migrate-up:
	@migrate -path migrations -database ${DB_URL} up

migrate-down:
	@migrate -path migrations -database ${DB_URL} down