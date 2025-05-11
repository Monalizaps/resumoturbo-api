# Makefile - comandos rápidos para o projeto ResumoTurbo

APP_NAME=resumoturbo-api

build:
	docker compose build

up:
	docker compose up

up-dev:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

restart:
	docker compose down && docker compose up --build

status:
	curl http://localhost:8080/status | jq

exec:
	docker exec -it $(APP_NAME) sh
