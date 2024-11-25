include .env

create_migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name using 'name=<migration_name>'"; \
		exit 1; \
	fi
	goose -dir=./database/migrations create $(name) sql

migrate_up:
	goose -dir=database/migrations postgres "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" up

migrate_down:
	goose -dir=database/migrations postgres "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" down

.PHONY: create_migration migrate_up migrate_down
