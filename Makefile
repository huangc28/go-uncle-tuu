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

PG_DSN=postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/darkpanda?sslmode=disable
PG_TEST_DSN=postgres://$(TEST_PG_USER):$(TEST_PG_PASSWORD)@$(TEST_PG_HOST):$(TEST_PG_PORT)/darkpanda?sslmode=disable

# Generate models from migration SQL schemas. This tool uses
# `https://github.com/kyleconroy/sqlc` to parse SQL syntax
# and generate corresponding models.
gen_model:
	gen-model gen --dbname=$(DB_NAME) --host=$(DB_HOST) --password=$(DB_PASSWORD) --port=$(DB_PORT) --username=$(DB_USER)

# Build & Deploy

# List of systemctl service name to host up app & worker.
APP_SERVICE_NAME                    = uncletuu.service
INVENTORY_IMPORTER_SERVICE_NAME     = uncletuu_inventory_importer.service

deploy:
	ssh -t $(DEPLOY_TARGET) 'cd /root/uncletuu/go-uncle-tuu && \
		git pull https://$(GITHUB_USER):$(GITHUB_ACCESS_TOKEN)@github.com/huangc28/go-uncle-tuu.git && \
		make build && \
		sudo systemctl stop $(APP_SERVICE_NAME) && \
		sudo systemctl start $(APP_SERVICE_NAME) && \
		sudo systemctl stop $(INVENTORY_IMPORTER_SERVICE_NAME) && \
		sudo systemctl start $(INVENTORY_IMPORTER_SERVICE_NAME)'

build: build_inventory_import_worker
	echo 'building production binary...'
	cd $(CURRENT_DIR)/cmd/app && GOOS=linux GOARCH=amd64 go build -o ../../bin/uncletuu_be -v .

build_inventory_import_worker:
	echo 'building inventory import worker'
	cd $(CURRENT_DIR)/cmd/inventory_import_worker && GOOS=linux GOARCH=amd64 go build -o ../../bin/inventory_importer_worker -v .

.PHONY: build
