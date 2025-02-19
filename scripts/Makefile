include ../.env

# Build the database URL from .env variables
DB_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

.PHONY: migrate-create migrate-up migrate-down migrate-force print-db-url

# Create a new migration file
migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

# Run all migrations
migrate-up:
	migrate -path ../internal/db/migrations -database "${DB_URL}" up

# Rollback all migrations
migrate-down:
	migrate -path ../internal/db/migrations -database "${DB_URL}" down

# Rollback one migration
migrate-down-1:
	migrate -path ../internal/db/migrations -database "${DB_URL}" down 1

# Force set version (useful if you get a dirty database error)
migrate-force:
	migrate -path ../internal/db/migrations -database "${DB_URL}" force $(version)

# Print the database URL (for debugging)
print-db-url:
	@echo ${DB_URL}
