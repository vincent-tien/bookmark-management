.PHONY: run swagger dev-run test

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

dev-run: swagger run

COVERAGE_THRESHOLD ?= 60
COVERAGE_OUT       ?= coverage.out
COVERAGE_HTML      ?= coverage.html

PKG_EXCLUDES := \
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

PKG_EXCLUDE_REGEX := $(subst $(space),|,$(strip $(foreach d,$(PKG_EXCLUDES),/$(d)($$|/))))

test:
	@set -eu; \
	pkgs=$$(go list ./... | grep -vE '$(PKG_EXCLUDE_REGEX)'); \
	coverpkgs=$$(echo "$$pkgs" | tr '\n' ',' | sed 's/,$$//'); \
	go test $$pkgs -covermode=atomic -coverprofile=$(COVERAGE_OUT) -coverpkg=$$coverpkgs -p 1; \
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML); \
	total=$$(go tool cover -func=$(COVERAGE_OUT) | awk '/^total:/ {gsub(/%/,"",$$3); print $$3}'); \
	awk -v t="$$total" -v th="$(COVERAGE_THRESHOLD)" 'BEGIN{ \
		if (t+0 < th+0) {printf "❌ Coverage (%.2f%%) is below threshold (%.2f%%)\n", t, th; exit 1} \
		else {printf "✅ Coverage (%.2f%%) meets threshold (%.2f%%)\n", t, th; exit 0} \
	}'
