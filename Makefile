.PHONY: up down logs lint

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f app

lint:
	golangci-lint run