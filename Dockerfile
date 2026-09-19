# Build the React app into ui/dist.
FROM node:22-alpine AS ui
WORKDIR /src/ui/app
COPY ui/app/package.json ui/app/package-lock.json ./
RUN npm ci
COPY ui/app/ ./
RUN npm run build -- --outDir /src/ui/dist --emptyOutDir

# Compile the Go binary with the built frontend embedded.
FROM golang:1.27-alpine AS backend
WORKDIR /src
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY --from=ui /src/ui/dist/ ./cmd/server/web/dist/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/cms ./cmd/server

# Runtime: one static binary, no database, no shell tooling.
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 cms
COPY --from=backend /out/cms /usr/local/bin/cms
USER cms
EXPOSE 8080
ENV CMS_HTTP_ADDR=:8080
# Mount a volume here and set CMS_CONTENT_FILE=/data/content.json to keep edits
# across restarts; without it the store is memory-only and the seed reloads.
VOLUME ["/data"]
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget -qO- http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/usr/local/bin/cms"]