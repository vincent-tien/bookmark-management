IMG_NAME := vincenttien/bookmark-service

GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)
BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)

# default
IMG_TAG := latest

# if main branch
ifeq ($(BRANCH),main)
  IMG_TAG := dev
endif

# if git tag exists → highest priority
ifneq ($(GIT_TAG),)
  IMG_TAG := $(GIT_TAG)
endif

export IMG_NAME IMG_TAG?=

.PHONY: run swagger dev-run test docker-test docker-lint docker-build docker-release docker-login migrate

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

dev-run: swagger run

COVERAGE_THRESHOLD ?= 50
COVERAGE_OUT       ?= coverage.out
COVERAGE_HTML      ?= coverage.html

COVERAGE_EXCLUDES := \
	mocks \
	docs \
	config \
	cmd \
	internal/routers \
	internal/errors \
	internal/dto \
	internal/test

empty :=
space := $(empty) $(empty)

COVERAGE_EXCLUDE_REGEX := $(subst $(space),|,$(strip $(foreach d,$(COVERAGE_EXCLUDES),/$(d)(/|:))))
COVERAGE_OUT      ?= coverage.out
COVERAGE_FILTERED ?= coverage.filtered.out
COVERAGE_HTML     ?= coverage.html

test:
	@set -eu; \
	go test ./... -covermode=atomic -coverprofile=$(COVERAGE_OUT) -coverpkg=./... -p 1; \
	if [ -n "$(COVERAGE_EXCLUDE_REGEX)" ]; then \
		grep -v -E "$(COVERAGE_EXCLUDE_REGEX)" $(COVERAGE_OUT) > $(COVERAGE_FILTERED) || cp $(COVERAGE_OUT) $(COVERAGE_FILTERED); \
	else \
		cp $(COVERAGE_OUT) $(COVERAGE_FILTERED); \
	fi; \
	go tool cover -html=$(COVERAGE_FILTERED) -o $(COVERAGE_HTML); \
	total=$$(go tool cover -func=$(COVERAGE_FILTERED) | awk '/^total:/ {gsub(/%/,"",$$3); print $$3}'); \
	awk -v t="$$total" -v th="$(COVERAGE_THRESHOLD)" 'BEGIN{ \
		if (t+0 < th+0) {printf "❌ Coverage (%.2f%%) is below threshold (%.2f%%)\n", t, th; exit 1} \
		else {printf "✅ Coverage (%.2f%%) meets threshold (%.2f%%)\n", t, th; exit 0} \
	}'

COVERAGE_FOLDER=./coverage

docker-test:
	@set -eu; \
	mkdir -p $(COVERAGE_FOLDER); \
	DOCKER_BUILDKIT=1 docker build \
		--target base \
		-t bookmark-service-test-base:dev \
		. ; \
	container_id=$$(docker run -d \
		-e COVERAGE_EXCLUDES="$(COVERAGE_EXCLUDE_REGEX)" \
		bookmark-service-test-base:dev \
		sh -ec '\
			mkdir -p /tmp/coverage && \
			go test ./... \
				-coverprofile=/tmp/coverage/coverage.tmp \
				-covermode=atomic \
				-coverpkg=./... \
				-p 1 \
				-timeout=10m && \
			if [ -z "$$COVERAGE_EXCLUDES" ]; then \
				cp /tmp/coverage/coverage.tmp /tmp/coverage/coverage.out; \
			else \
				grep -v -E "$$COVERAGE_EXCLUDES" /tmp/coverage/coverage.tmp > /tmp/coverage/coverage.out || cp /tmp/coverage/coverage.tmp /tmp/coverage/coverage.out; \
			fi && \
			go tool cover -html=/tmp/coverage/coverage.out -o /tmp/coverage/coverage.html \
		'); \
	exit_code=$$(docker wait $$container_id); \
	if [ $$exit_code -ne 0 ]; then \
		echo "Test execution failed. Container logs:"; \
		docker logs $$container_id; \
		docker rm $$container_id; \
		exit $$exit_code; \
	fi; \
	docker cp $$container_id:/tmp/coverage/coverage.out $(COVERAGE_FOLDER)/coverage.out; \
	docker cp $$container_id:/tmp/coverage/coverage.html $(COVERAGE_FOLDER)/coverage.html; \
	docker rm $$container_id; \
	if [ ! -f "$(COVERAGE_FOLDER)/coverage.out" ]; then \
		echo "❌ coverage.out not found in $(COVERAGE_FOLDER)"; \
		exit 1; \
	fi; \
	total=$$(go tool cover -func="$(COVERAGE_FOLDER)/coverage.out" | awk '/^total:/ {gsub(/%/,"",$$3); print $$3}'); \
	awk -v t="$$total" -v th="$(COVERAGE_THRESHOLD)" 'BEGIN{ \
		if (t+0 < th+0) {printf "❌ Coverage (%.2f%%) is below threshold (%.2f%%)\n", t, th; exit 1} \
		else {printf "✅ Coverage (%.2f%%) meets threshold (%.2f%%)\n", t, th; exit 0} \
	}'

GOLANG_LINT_VERSION ?= v2.7.2

docker-lint:
	@set -eu; \
	docker buildx build \
		--build-arg GOLANG_LINT_VERSION="$(GOLANG_LINT_VERSION)" \
		--target golangci-lint \
		-t bookmark-service-lint:dev \
		--progress=plain \
		.

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

docker-release: docker-build
	docker push $(IMG_NAME):$(IMG_TAG)

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

# Migration commands using goose
# Database connection variables (can be overridden via environment)
DB_HOST     ?= localhost
DB_USER     ?= ebvn
DB_PASSWORD ?= abc123
DB_NAME     ?= ebvn_bm
DB_PORT     ?= 5432

# Goose configuration
GOOSE_MIGRATION_DIR ?= migrations
GOOSE_DRIVER        ?= postgres
GOOSE_DBSTRING      ?= host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) dbname=$(DB_NAME) password=$(DB_PASSWORD) sslmode=disable

# Get the command from arguments (everything after 'migrate')
MIGRATE_CMD := $(word 2,$(MAKECMDGOALS))
# For create command, get the migration name and type (default type is sql)
MIGRATE_NAME := $(word 3,$(MAKECMDGOALS))
MIGRATE_TYPE := $(or $(word 4,$(MAKECMDGOALS)),sql)

migrate:
	@if [ -z "$(MIGRATE_CMD)" ]; then \
		echo "Error: Command is required. Usage: make migrate <command>"; \
		echo "Available commands: up, down, up-to, down-to, status, version, create"; \
		exit 1; \
	fi
	@case "$(MIGRATE_CMD)" in \
		up) \
			goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up \
			;; \
		down) \
			goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down \
			;; \
		up-to) \
			if [ -z "$(VERSION)" ]; then \
				echo "Error: VERSION is required. Usage: make migrate up-to VERSION=<version>"; \
				exit 1; \
			fi; \
			goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up-to $(VERSION) \
			;; \
		down-to) \
			if [ -z "$(VERSION)" ]; then \
				echo "Error: VERSION is required. Usage: make migrate down-to VERSION=<version>"; \
				exit 1; \
			fi; \
			goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down-to $(VERSION) \
			;; \
		status) \
			goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" status \
			;; \
		version) \
			goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" version \
			;; \
		create) \
			if [ -z "$(MIGRATE_NAME)" ]; then \
				echo "Error: Migration name is required. Usage: make migrate create <migration_name> [sql|go]"; \
				exit 1; \
			fi; \
			goose -dir $(GOOSE_MIGRATION_DIR) create $(MIGRATE_NAME) $(MIGRATE_TYPE) \
			;; \
		*) \
			echo "Error: Unknown command '$(MIGRATE_CMD)'"; \
			echo "Available commands: up, down, up-to, down-to, status, version, create"; \
			exit 1 \
			;; \
	esac

# Prevent Make from trying to build migration command arguments as targets
%:
	@:

