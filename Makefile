PROJECT_NAME := product-reservation-system
MIGRATOR_DIR := ./cmd/migrator
MIGRATIONS_DIR := ./migrations
APP_DIR := ./cmd
STORAGE_DSN := "host=localhost port=5433 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"

GO := go
GO_BUILD := $(GO) build
GO_RUN := $(GO) run
GO_TEST := $(GO) test -v ./...
GO_MOD_TIDY := &(GO) mod tidy

MIGRATE_UP := $(GO_RUN) $(MIGRATOR_DIR) -command=up 
MIGRATE_DOWN := $(GO_RUN) $(MIGRATOR_DIR) -command=down
MIGRATE_STATUS := $(GO_RUN) $(MIGRATOR_DIR) -command=status

all: tidy build

# Build the project
build:
	@echo "Building $(PROJECT_NAME)..."
	$(GO_BUILD) -o bin/$(PROJECT_NAME) $(APP_DIR)

# Run the application
run:
	@echo "Starting application..."
	$(GO_RUN) $(APP_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GO_TEST)

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO_MOD_TIDY)

# Migrations
migrate-up:
	@echo "Applying migrations..."
	$(MIGRATE_UP)

migrate-down:
	@echo "Reverting migrations..."
	$(MIGRATE_DOWN)

migrate-status:
	@echo "Migration status:"
	$(MIGRATE_STATUS)

