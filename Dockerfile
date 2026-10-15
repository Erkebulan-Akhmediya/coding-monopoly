# syntax=docker/dockerfile:1

# ---- 1. Build the Vue client ------------------------------------------------
FROM node:22-alpine AS client
WORKDIR /src/client

COPY client/package.json client/package-lock.json ./
RUN npm ci

COPY client/ ./
# VITE_WS_BASE_URL must never be baked into a production build (see Makefile).
RUN env -u VITE_WS_BASE_URL npm run build && test -f dist/index.html
RUN if grep -rnE 'localhost:8080|ws://localhost|wss://localhost' dist/assets; then \
        echo "embedded JS contains localhost WS URL — refuse to ship" >&2; exit 1; \
    fi

# ---- 2. Build the Go binaries (frontend embedded via go:embed) --------------
FROM golang:1.26-alpine AS server
WORKDIR /src/server

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./
RUN rm -rf cmd/server/dist && mkdir -p cmd/server/dist
COPY --from=client /src/client/dist/ cmd/server/dist/

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/monopoly-server ./cmd/server \
    && go build -trimpath -ldflags="-s -w" -o /out/monopoly-seed ./cmd/seed

# ---- 3. Minimal runtime -----------------------------------------------------
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app

COPY --from=server /out/monopoly-server /out/monopoly-seed /usr/local/bin/

USER 10001:10001
WORKDIR /app
ENV LISTEN_ADDR=0.0.0.0:8080
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/monopoly-server"]
