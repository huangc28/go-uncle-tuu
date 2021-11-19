-include ./go-migrate-makefile/go_migrate_makefile

CURRENT_DIR = $(shell pwd)

# Export variables in .env
ifneq (,$(wildcard ./.env))
	include ./.env
	export
endif

run_local: run_local_docker
	go run cmd/main.go

run_local_docker:
	docker-compose \
		-f ./docker-compose.yaml \
		--env-file .env up \
		-d
