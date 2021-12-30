-include ./go-migrate-makefile/go_migrate.mk

CURRENT_DIR = $(shell pwd)

# Export variables in .env
ifneq (,$(wildcard ./.env))
	include ./.env
	export
endif

run_local: run_local_docker
	go mod tidy && go run cmd/main.go

run_local_docker:
	docker-compose \
		-f ./docker-compose.yaml \
		--env-file .env up \
		-d

# Build & Deploy

# List of systemctl service name to host up worker.
APP_SERVICE_NAME                    = uncletuu.service


deploy: build
	ssh -t root@api.darkpanda.love 'cd /root/uncletuu/go-uncle-tuu && \
		git pull https://$(GITHUB_USER):$(GITHUB_ACCESS_TOKEN)@github.com/huangc28/go-uncle-tuu.git && \
		make build && \
		sudo systemctl stop $(APP_SERVICE_NAME) && \
		sudo systemctl start $(APP_SERVICE_NAME)'

build:
	echo 'building production binary...'
	cd $(CURRENT_DIR)/cmd/app && GOOS=linux GOARCH=amd64 go build -o ../../bin/uncletuu_be -v .

.PHONY: build
