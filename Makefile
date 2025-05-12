# Database connection string
DB_STRING = postgres://howl:turnip_man1234@localhost:5432/social?sslmode=disable

# Migration commands
.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-version

# Apply all pending migrations
migrate-up:
	@echo "Applying migrations..."
	goose -dir migrations postgres "$(DB_STRING)" up

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	goose -dir migrations postgres "$(DB_STRING)" down

# Show migration status
migrate-status:
	@echo "Migration status:"
	goose -dir migrations postgres "$(DB_STRING)" status

# Create a new migration file
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name required"; \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)"
	goose -dir migrations create "$(name)" sql

# Show current database version
migrate-version:
	@echo "Current database version:"
	goose -dir migrations postgres "$(DB_STRING)" version

# Help command
help:
	@echo "Available commands:"
	@echo "  make migrate-up              - Apply all pending migrations"
	@echo "  make migrate-down            - Rollback the last migration"
	@echo "  make migrate-status          - Show migration status"
	@echo "  make migrate-create name=xxx - Create a new migration file"
	@echo "  make migrate-version         - Show current database version"
	@echo "  make help                    - Show this help message" 