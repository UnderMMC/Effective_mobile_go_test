MIGRATIONS_DIR=migrations
DB_URL=postgres://postgres:pgpwd4habr@localhost/postgres?sslmode=disable

.PHONY: migrate up down

migrate:
	@echo "Запуск миграций базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

rollback:
	@echo "Откат последней миграции базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

