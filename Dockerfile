ARG NODE_VERSION=25
ARG GO_VERSION=1.25.4

# Build stage for Node.js assets
FROM node:${NODE_VERSION}-alpine AS node-builder

WORKDIR /app

# Copy package files
COPY assets/package.json assets/pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm i -g pnpm
RUN pnpm install --frozen-lockfile

# Copy CSS assets and build
COPY assets/ ./
COPY tmpl ./tmpl
RUN pnpm run build

# Build stage for Go application
FROM golang:${GO_VERSION}-alpine AS go-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server -tags PROD .

# Final stage
FROM alpine:latest

ARG GOOSE_VERSION=3.26.0
ENV GOOSE_VERSION=$GOOSE_VERSION

RUN apk --no-cache add ca-certificates tzdata curl

RUN /bin/sh -c 'set -ex && \
    ARCH=$(uname -m); \
    case "$ARCH" in \
      aarch64) ARCH=arm64 ;; \
      *) ;; \
    esac; \
    curl --output /usr/local/bin/goose -LJO https://github.com/pressly/goose/releases/download/v$GOOSE_VERSION/goose_linux_$ARCH; \
    chmod +x /usr/local/bin/goose'

WORKDIR /app

COPY --from=go-builder /app/server ./
COPY --from=go-builder /app/i18n/*.json ./i18n/
COPY --from=node-builder /app/dist ./assets/dist
COPY --from=go-builder /app/db/migrations ./db/migrations
COPY --from=go-builder /app/scripts/entrypoint.sh ./

# Expose port
EXPOSE 3000

# Run the application
CMD ["/app/entrypoint.sh"]
