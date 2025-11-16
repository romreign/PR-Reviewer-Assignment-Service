.PHONY: help build run test test-coverage test-short lint fmt tidy \
        docker-build docker-up docker-down docker-logs docker-clean \
        clean dev ci migrate-up

UNAME_S := $(shell uname -s)
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
else
    DETECTED_OS := $(UNAME_S)
endif

BINARY_NAME := pr-reviewer-app
BINARY_EXT :=
ifeq ($(DETECTED_OS),Windows)
    BINARY_EXT := .exe
    MKDIR := mkdir
    RM := rmdir /s /q
    RM_FILE := del /q
    SHELL := cmd.exe
else
    MKDIR := mkdir -p
    RM := rm -rf
    RM_FILE := rm -f
endif

DOCKER_IMAGE := pr-reviewer:latest
GO := go
GOFLAGS := -v
BIN_PATH := bin/$(BINARY_NAME)$(BINARY_EXT)

help:
	@echo "PR Reviewer Assignment Service - Make Targets"
	@echo "=============================================="
	@echo ""
	@echo "Build & Run:"
	@echo "  make build           Build the application"
	@echo "  make run             Build and run locally"
	@echo "  make dev             Run with auto-reload (requires air)"
	@echo ""
	@echo "Testing:"
	@echo "  make test            Run all tests"
	@echo "  make test-coverage   Run tests with coverage report"
	@echo "  make test-short      Run tests (short mode)"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt             Format code and organize imports"
	@echo "  make lint            Run linter (golangci-lint)"
	@echo "  make tidy            Tidy dependencies"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build    Build Docker image"
	@echo "  make docker-up       Start services (docker-compose)"
	@echo "  make docker-down     Stop services"
	@echo "  make docker-logs     Show service logs"
	@echo "  make docker-clean    Stop and remove containers/images"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up      Run database migrations"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean           Clean build artifacts"
	@echo "  make ci              Run CI pipeline (lint + test)"
	@echo ""

build:
	@echo "[BUILD] Compiling $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BIN_PATH) ./cmd/api
	@echo "[OK] Build complete: $(BIN_PATH)"

run: build
	@echo "[RUN] Starting application..."
ifeq ($(DETECTED_OS),Windows)
	$(BIN_PATH)
else
	./$(BIN_PATH)
endif

test:
	@echo "[TEST] Running all tests..."
	$(GO) test $(GOFLAGS) ./...

test-coverage:
	@echo "[TEST] Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "[OK] Coverage report generated: coverage.html"

test-short:
	@echo "[TEST] Running short tests..."
	$(GO) test -short ./...

lint:
	@echo "[LINT] Running golangci-lint..."
ifeq ($(DETECTED_OS),Windows)
	@where golangci-lint >nul 2>&1 || (echo [INFO] Installing golangci-lint... && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
else
	@which golangci-lint >/dev/null || (echo "[INFO] Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
endif
	golangci-lint run ./...

fmt:
	@echo "[FMT] Formatting code..."
	$(GO) fmt ./...
ifeq ($(DETECTED_OS),Windows)
	@where goimports >nul 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
else
	@which goimports >/dev/null || go install golang.org/x/tools/cmd/goimports@latest
endif
	goimports -w .
	@echo "[OK] Code formatted"

tidy:
	@echo "[TIDY] Managing dependencies..."
	$(GO) mod tidy
	$(GO) mod verify
	@echo "[OK] Dependencies tidied"

docker-build:
	@echo "[DOCKER] Building image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "[OK] Image built successfully"

docker-up:
	@echo "[DOCKER] Starting services..."
	docker-compose up -d --build
	@echo "[OK] Services started"
	docker-compose logs -f

docker-down:
	@echo "[DOCKER] Stopping services..."
	docker-compose down
	@echo "[OK] Services stopped"

docker-logs:
	@echo "[DOCKER] Showing logs..."
	docker-compose logs -f

docker-clean:
	@echo "[DOCKER] Cleaning up..."
	docker-compose down -v
ifeq ($(DETECTED_OS),Windows)
	@docker rmi $(DOCKER_IMAGE) 2>nul || @echo "[INFO] Image cleanup skipped"
else
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || echo "[INFO] Image cleanup skipped"
endif
	@echo "[OK] Docker cleaned"

clean:
	@echo "[CLEAN] Removing build artifacts..."
ifeq ($(DETECTED_OS),Windows)
	@if exist bin $(RM) bin
	@$(RM_FILE) coverage.out 2>nul || @echo
	@$(RM_FILE) coverage.html 2>nul || @echo
else
	$(RM) bin/
	$(RM_FILE) coverage.out coverage.html
endif
	$(GO) clean
	@echo "[OK] Cleaned"

dev:
	@echo "[DEV] Starting with auto-reload..."
ifeq ($(DETECTED_OS),Windows)
	@where air >nul 2>&1 || (echo [INFO] Installing air... && go install github.com/cosmtrek/air@latest)
else
	@which air >/dev/null || (echo "[INFO] Installing air..." && go install github.com/cosmtrek/air@latest)
endif
	air

migrate-up:
	@echo "[MIGRATE] Running migrations..."
	psql -U postgres -d pr_review_db -f migrations/001_init.sql
	@echo "[OK] Migrations complete"

ci: lint test
	@echo "[CI] Pipeline complete"

