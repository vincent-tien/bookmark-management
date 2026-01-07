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

.PHONY: run swagger dev-run test docker-test docker-lint docker-build docker-release docker-login

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
	docker buildx build \
		--build-arg COVERAGE_EXCLUDES="$(COVERAGE_EXCLUDE_REGEX)" \
		--target test-coverage \
		-t bookmark-service-test:dev \
		--output type=local,dest=$(COVERAGE_FOLDER) . ; \
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
