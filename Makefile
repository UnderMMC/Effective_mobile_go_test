MIGRATIONS_DIR=migrations
DB_URL=postgres://postgres:pgpwd4habr@localhost/postgres?sslmode=disable

.PHONY: migrateup migratedown

migrateup:
	@echo "Запуск миграций базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migratedown:
	@echo "Откат последней миграции базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

run:
	go run cmd/sum/main.go