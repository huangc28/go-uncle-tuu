-include ./go-migrate-makefile/go_migrate.mk

CURRENT_DIR = $(shell pwd)

# Export variables in .env
ifneq (,$(wildcard ./.env))
	include ./.env
	export
endif

ifeq (, $(shell which gen-model))
	$(error "No gen-model in $(GOPATH)/bin, please install git@github.com:huangc28/go-migration-model-generator.git before proceeding")
endif

run_local: run_local_docker
	go mod tidy && go run cmd/app/main.go

run_local_docker:
	docker-compose \
		-f ./docker-compose.yaml \
		--env-file .env up \
		-d

MIGRATE_CMD=migrate
MIGRATE_CREATE_CMD=create
MIGRATE_UP_CMD=up
MIGRATE_DOWN_CMD=down

PG_DSN=postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/darkpanda?sslmode=disable
PG_TEST_DSN=postgres://$(TEST_PG_USER):$(TEST_PG_PASSWORD)@$(TEST_PG_HOST):$(TEST_PG_PORT)/darkpanda?sslmode=disable

# Generate models from migration SQL schemas. This tool uses
# `https://github.com/kyleconroy/sqlc` to parse SQL syntax
# and generate corresponding models.
gen_model:
	gen-model gen --dbname=$(DB_NAME) --host=$(DB_HOST) --password=$(DB_PASSWORD) --port=$(DB_PORT) --username=$(DB_USER)

# Build & Deploy

# List of systemctl service name to host up worker.
APP_SERVICE_NAME                    = uncletuu.service

deploy: build
	ssh -t $(DEPLOY_TARGET) 'cd /root/uncletuu/go-uncle-tuu && \
		git pull https://$(GITHUB_USER):$(GITHUB_ACCESS_TOKEN)@github.com/huangc28/go-uncle-tuu.git && \
		make build && \
		sudo systemctl stop $(APP_SERVICE_NAME) && \
		sudo systemctl start $(APP_SERVICE_NAME)'

build:
	echo 'building production binary...'
	cd $(CURRENT_DIR)/cmd/app && GOOS=linux GOARCH=amd64 go build -o ../../bin/uncletuu_be -v .

.PHONY: build
