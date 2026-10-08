include .env
export

export PROJECT_ROOT=$(shell pwd)

run:
	@go run cmd/financer/main.go

env-up:
	@docker compose up -d financer_postgres

env-cleanup:
	@docker compose down -v financer_postgres

env-port-forward:
	@docker compose up -d port_forwarder

env-port-close:
	@docker compose down port_forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "'seq' parameter is required. Try 'make migrate-create seq=init'"; \
		exit 1; \
	fi; \
	docker compose run --rm financer_postgres_migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "'action' parameter is required. Try 'make migrate-action action=up'"; \
		exit 1; \
	fi; \
	docker compose run --rm financer_postgres_migrate \
		-path migrations -database postgres://${DB_USER}:${DB_PASSWORD}@financer_postgres:${DB_PORT}/${DB_NAME}?sslmode=disable \
		"$(action)"

docker-ps:
	@docker compose ps

financer-deploy:
	@docker compose up -d --build financer

financer-undeploy:
	@docker compose down