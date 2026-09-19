SHELL := /bin/sh

COMPOSE ?= docker compose
GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose@v3.24.3
POSTGRES_URL ?= postgres://cms:cms@localhost:5432/cms?sslmode=disable

.PHONY: up down migrate seed backend frontend dev

up:
	$(COMPOSE) up -d --wait

down:
	$(COMPOSE) down

migrate:
	$(GOOSE) -dir cmd/server/migrations/postgres postgres "$(POSTGRES_URL)" up

seed:
	psql "$(POSTGRES_URL)" -v ON_ERROR_STOP=1 -f scripts/seed_content.sql

backend:
	set -a; [ ! -f .env ] || . ./.env; set +a; \
	CMS_APPLICATION_DATABASE_URL="$${CMS_APPLICATION_DATABASE_URL:-host=localhost user=cms password=cms dbname=cms port=5432 sslmode=disable}" \
	go run ./cmd/server

frontend:
	@if [ ! -x ui/app/node_modules/.bin/vite ]; then \
		npm --prefix ui/app ci; \
	fi
	npm --prefix ui/app run dev

dev: up migrate
	@trap 'kill 0' INT TERM; \
		$(MAKE) backend & \
		$(MAKE) frontend & \
		wait
