FROM node:22-alpine AS web
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.26-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
RUN mkdir -p /data/uploads
COPY --from=backend /server /app/server
COPY --from=backend /app/backend/migrations /app/migrations
COPY --from=web /app/web/dist /app/web/dist
ENV PORT=8080
ENV DB_PATH=/data/speedmatch.db
ENV UPLOAD_DIR=/data/uploads
ENV STATIC_DIR=/app/web/dist
ENV CORS_ORIGIN=*
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1
CMD ["/app/server"]
