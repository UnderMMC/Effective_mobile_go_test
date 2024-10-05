MIGRATIONS_DIR=migrations
DB_URL=postgres://postgres:pgpwd4habr@localhost:5432/postgres?sslmode=disable

.PHONY: up down run

up:
	@echo "Запуск миграций базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

down:
	@echo "Откат последней миграции базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

run:
	@echo "Применение миграций базы данных..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
	@echo "Запуск сервиса..."
	go run cmd/sum/main.go