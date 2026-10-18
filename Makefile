include .env
export

export PROJECT_ROOT=$(shell pwd)

COMPOSE ?= docker compose

.PHONY: run env-up env-down env-cleanup env-port-forward env-port-close migrate-create migrate-up migrate-down migrate-action psql logs docker-ps swagger-gen financer-deploy financer-undeploy redis-up redis-down redis-restart redis-logs redis-cli redis-ping

run:
	@go run cmd/financer/main.go

env-up:
	@$(COMPOSE) up -d financer_postgres

env-down:
	@$(COMPOSE) down financer_postgres

env-cleanup:
	@$(COMPOSE) down -v financer_postgres

env-port-forward:
	@$(COMPOSE) up -d port_forwarder

env-port-close:
	@$(COMPOSE) down port_forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "'seq' parameter is required. Try 'make migrate-create seq=init'"; \
		exit 1; \
	fi; \
	$(COMPOSE) run --rm financer_postgres_migrate \
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
	$(COMPOSE) run --rm financer_postgres_migrate \
		-path migrations -database postgres://${DB_USER}:${DB_PASSWORD}@financer_postgres:${DB_PORT}/${DB_NAME}?sslmode=disable \
		$(action)

psql:
	@docker exec -it financer_postgres psql -U financer -d financer

logs:
	@docker logs financer

docker-ps:
	@$(COMPOSE) ps

swagger-gen:
	@$(COMPOSE) run --rm swagger \
		init \
		-g cmd/financer/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

financer-deploy:
	@$(COMPOSE) up -d --build financer

financer-undeploy:
	@$(COMPOSE) down

redis-up:
	@$(COMPOSE) up -d redis

redis-down:
	@$(COMPOSE) stop redis

redis-restart:
	@$(COMPOSE) restart redis

redis-logs:
	@$(COMPOSE) logs -f --tail=100 redis

redis-cli:
	@$(COMPOSE) exec redis redis-cli -a ${REDIS_PASSWORD}

redis-ping:
	@$(COMPOSE) exec redis redis-cli -a ${REDIS_PASSWORD} ping