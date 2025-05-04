.PHONY: up down psql run

# Start PostgreSQL container
up:
	docker compose up -d

# Stop PostgreSQL container
down:
	docker compose down

# Connect to PostgreSQL
psql:
	docker exec -it crm_postgres psql -U postgres -d crm_db

# Run the application
run:
	go run cmd/backend_crm/main.go

# Clean up volumes
clean:
	docker compose down -v


