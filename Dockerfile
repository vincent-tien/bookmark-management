# =====================
# Stage : Base
# =====================
FROM golang:alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

# Install dependencies needed for CGO and SQLite
# gcc, musl-dev, and sqlite-dev are required for go-sqlite3 to work
RUN apk add --no-cache ca-certificates gcc musl-dev sqlite-dev

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# =====================
# Stage : Build binary
# =====================

FROM base AS builder

# Build binary, tắt CGO để chạy trên alpine
# Build binary static + strip debug
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o bookmark_service \
    cmd/api/main.go


 # =====================
 # Stage : Test exec
 # =====================
 # Note: This stage is kept for compatibility, but tests are now run
 # in a container using docker run (see Makefile docker-test target)
 # to avoid network binding restrictions during docker build

FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDES=""

ENV GO_TEST_TIMEOUT=10m

RUN sh -ec '\
  mkdir -p "${_outputdir}" && \
  go test ./... \
    -coverprofile=coverage.tmp \
    -covermode=atomic \
    -coverpkg=./... \
    -p 1 \
    -timeout=10m && \
  if [ -z "${COVERAGE_EXCLUDES}" ]; then \
    cp coverage.tmp "${_outputdir}/coverage.out"; \
  else \
    grep -v -E "${COVERAGE_EXCLUDES}" coverage.tmp > "${_outputdir}/coverage.out" || cp coverage.tmp "${_outputdir}/coverage.out"; \
  fi && \
  go tool cover -html="${_outputdir}/coverage.out" -o "${_outputdir}/coverage.html" \
'

FROM scratch AS test-coverage
ARG _outputdir="/tmp/coverage"

COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /


# =====================
# Stage: Lint
# =====================
FROM base AS golangci-lint

ARG GOLANG_LINT_VERSION=v1.56.2

RUN apk add --no-cache curl git

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b /usr/local/bin ${GOLANG_LINT_VERSION}

RUN golangci-lint run


# =====================
# Stage 2: Runtime
# =====================
FROM alpine AS final

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=builder /opt/app/bookmark_service .
COPY --from=builder /opt/app/docs .

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN chown -R nonroot:nonroot /app

USER nonroot:nonroot

CMD ["/app/bookmark_service"]