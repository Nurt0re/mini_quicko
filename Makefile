.PHONY: build up stop migrate

build: ## Build Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d

stop: ## Stop all services
	docker-compose stop

migrate: ## Run database migrations
	docker-compose exec postgres psql -U admin -d mini_quicko -f /docker-entrypoint-initdb.d/init.sql

