.PHONY: dev backend frontend build migrate seed docker

dev:
	@echo "Run in two terminals:"
	@echo "  make backend"
	@echo "  make frontend"

backend:
	cd backend && go run ./cmd/server

frontend:
	cd web && npm run dev

build:
	cd web && npm run build
	cd backend && go build -o bin/server ./cmd/server

migrate:
	cd backend && go run ./cmd/server

seed:
	cd backend && go run ./cmd/seed

docker:
	docker compose up --build
