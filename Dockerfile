# =====================
# Stage 1: Build
# =====================
FROM golang:alpine AS builder

RUN mkdir -p /opt/app

WORKDIR /opt/app

# Chỉ cài thứ cần thiết
RUN apk add --no-cache ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary, tắt CGO để chạy trên alpine
# Build binary static + strip debug
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o bookmark_service \
    cmd/api/main.go

# =====================
# Stage 2: Runtime
# =====================
FROM alpine AS run

RUN apk add --no-cache ca-certificates \
    && addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

WORKDIR /app

COPY --from=builder /opt/app/bookmark_service .
COPY --from=builder /opt/app/docs .

USER nonroot:nonroot

CMD ["/app/bookmark_service"]