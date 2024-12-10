# Variables
APP_CMD = cmd/main.go
APP_NAME = coffee-shop

# Tasks
.PHONY: all run lint start

# Default target
run: start

# Target to run the application
start:
	@echo "Starting the application..."
	@go run $(APP_CMD)

# Target to run docker
docker-up:
	@echo "Starting docker"
	@docker-compose up mysql redis kafka mailhog -d
	@echo "All docker is up"

# Target to run docker
docker-down:
	@echo "Stopping docker"
	@docker-compose down mysql redis kafka mailhog
	@echo "All docker is down"

# Help target to display usage
help:
	@echo "Usage:"
	@echo "  make run   - Run the application with syntax check and tests"
	@echo "  make docker-up  - Run docker for mysql redis kafka mailhog"
	@echo "  make docker-down  - Stopping docker for mysql redis kafka mailhog"