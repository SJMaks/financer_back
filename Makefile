include .env
export

run:
	@go run cmd/financer/main.go

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down