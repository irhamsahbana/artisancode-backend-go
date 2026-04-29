.PHONY: migrate restore create-migration seed dev daemon ws build build-dev build-staging build-production lint-fix lint-ci nats commit help mock test-storage-upload cleanup-storage-orphans cleanup-message-queue scheduler

# Variables
GOOSE_CMD := goose
APP_NAME := app
APP_PORT := 3939
MOCKERY_VERSION := v3.7.0

# Default
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Database Migration
migrate: ## Run migrations (usage: make migrate cmd=up)
	@$(GOOSE_CMD) $(cmd)

restore: ## Restore PostgreSQL from dump (usage: make restore file=backup.sql)
	@PGPASSWORD=$(POSTGRES_PASSWORD) pg_restore --no-owner --clean -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) ./backups/$(file)

create-migration: ## Create migration (usage: make create-migration name=create_users_table)
	@$(GOOSE_CMD) create $(name) sql

seed: ## Run seed (usage: make seed table=rbac)
	@go run ./cmd/bin/main.go seed -table=$(table)

# Development
dev: ## Run development server
	@go run ./cmd/bin/main.go --port=$(APP_PORT)

mock: ## Generate mocks using mockery
	@go run github.com/vektra/mockery/v3@$(MOCKERY_VERSION)

test-storage-upload: ## Test storage presigned upload (usage: AUTH_TOKEN=xxx make test-storage-upload [file=./path.jpg])
	@./scripts/test_storage_upload.sh $(file)

cleanup-storage-orphans: ## Clean expired pending storage files (usage: make cleanup-storage-orphans [limit=100])
	@go run ./cmd/bin/main.go cronjob --task=cleanup-expired-storage-files --limit=$(or $(limit),100)

cleanup-message-queue: ## Clean processed message queue rows (usage: make cleanup-message-queue [limit=500] [retention_hours=168])
	@go run ./cmd/bin/main.go cronjob --task=cleanup-processed-message-queue --limit=$(or $(limit),500) --retention-hours=$(or $(retention_hours),168)

scheduler: ## Run internal cron scheduler
	@go run ./cmd/bin/main.go scheduler

daemon: ## Run with daemon (pmgo)
	@pmgo

ws: ## Run websocket server on port 8080
	@go run ./cmd/bin/main.go ws --port=8080

# Build
build: ## Build the application
	@go build -o ./$(APP_NAME) ./cmd/bin/main.go

build-dev: ## Build and deploy to dev
	@git pull
	@go build -o ./$(APP_NAME) ./cmd/bin/main.go
	@immortalctl stop kpf-dev
	@mv ./$(APP_NAME) ../binaries/kpf-dev
	@immortalctl start kpf-dev
	@immortalctl status

build-staging: ## Build and deploy to staging
	@git pull
	@go build -o ./$(APP_NAME) ./cmd/bin/main.go
	@immortalctl stop kpf-staging
	@mv ./$(APP_NAME) ../binaries/kpf-staging
	@immortalctl start kpf-staging
	@immortalctl status

build-production: ## Build and deploy to production
	@git pull
	@go build -o ./$(APP_NAME) ./cmd/bin/main.go
	@immortalctl stop kpf-production
	@mv ./$(APP_NAME) ../binaries/kpf-production
	@immortalctl start kpf-production
	@immortalctl status

# Linting
lint-fix: ## Auto-fix linting issues
	@gofmt -w .

lint-ci: ## Run linting for CI
	@golangci-lint run

# Infrastructure
nats: ## Start NATS server with JetStream
	@nats-server --js

# Git
commit: ## Commit with linting (usage: make commit msg="Your commit message")
	@task lint-fix
	@task lint-ci
	@git add .
	@git commit -m "$(msg)"
