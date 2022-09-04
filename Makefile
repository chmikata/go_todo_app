.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## build docker image to deploy
	docker build -t chmikata/gotodo:${DOCKER_TAG} --target deploy ./

build-local: ## Build docker image to local development
	docker compose build --no-cache

up: ## Do docker compose up with hot reload
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

test: ## Execute tests
	go test -race -shuffle=on ./...

migrate: ## Migration DB
	psqldef -U todo -W todo -h localhost -p 5432 todo -f ./_tools/postgres/schema.sql

dry-migrate: ## Dry-run Migration DB
	psqldef -U todo -W todo -h localhost -p 5432 todo --dry-run -f ./_tools/postgres/schema.sql

generate: ## Generate codes
	go generate ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
