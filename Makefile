.PHONY: run database seed docker-up docker-down

# Main command
app: run

# Subcommands
run:
	@echo "Running the application..."
	go run main.go run

# Create database tables
database:
	@echo "Creating database tables..."
	go run main.go create database

# Seed the database
seed:
	@echo "Creating database tables if not exists..."
	go run main.go create database
	@echo "Seeding the database with mocked data..."
	go run main.go create seed

# Clear database data
clear:
	@echo "Delete all database data..."
	go run main.go delete tables

# Docker commands
docker-up: 
	@echo "Starting docker containers..."
	docker compose -f deployments/docker-compose-dev.yaml up -d

docker-down:
	@echo "Removing docker containers and volumes..."
	docker compose -f deployments/docker-compose-dev.yaml down -v

coverage:
	@echo "Running test coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

ocov:
	open coverage.html

