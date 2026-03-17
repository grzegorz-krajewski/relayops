COMPOSE_FILE=deploy/compose/docker-compose.yml

up:
	docker compose --env-file .env -f $(COMPOSE_FILE) up -d --build

down:
	docker compose --env-file .env -f $(COMPOSE_FILE) down

logs:
	docker compose --env-file .env -f $(COMPOSE_FILE) logs -f

ps:
	docker compose --env-file .env -f $(COMPOSE_FILE) ps

restart:
	docker compose --env-file .env -f $(COMPOSE_FILE) restart

build:
	docker compose --env-file .env -f $(COMPOSE_FILE) build

test:
	@echo "Tests will be added in next sprint"

lint:
	@echo "Linters will be added in next sprint"

proto:
	@echo "Proto generation will be added in next sprint"