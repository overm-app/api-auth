include .env
export

MIGRATE=~/go/bin/migrate
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate-up:
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" down 1

migrate-version:
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" version

migrate-create:
	$(MIGRATE) create -ext sql -dir ./migrations -seq $(name)

db-reset:
	docker compose down -v postgres
	docker compose up -d postgres
	@sleep 5
	$(MAKE) migrate-up