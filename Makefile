.PHONY: help dev frontend backend build clean install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Run in development mode (frontend + backend)
	@echo "Starting development environment..."
	@echo "Frontend: http://localhost:5173"
	@echo "Backend: http://localhost:8080"
	@cd frontend && npm run dev & \
	go run main.go

frontend: ## Build frontend only
	@echo "Building frontend..."
	@cd frontend && npm install && npm run build

backend: ## Build backend only
	@echo "Building backend..."
	@go build -o supervisord-monitor

build: ## Build complete application (frontend + backend)
	@echo "Building complete application..."
	@$(MAKE) frontend
	@$(MAKE) backend

windows: ## Build Windows executable
	@echo "Building Windows executable..."
	@cd frontend && npm install && npm run build
	@set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=1&& \
		go build -ldflags "-s -w" -o supervisord-monitor.exe

linux: ## Build Linux executable
	@echo "Building Linux executable..."
	@cd frontend && npm install && npm run build
	@export GOOS=linux && export GOARCH=amd64 && export CGO_ENABLED=0 && \
		go build -ldflags "-s -w" -o supervisord-monitor

darwin: ## Build macOS executable
	@echo "Building macOS executable..."
	@cd frontend && npm install && npm run build
	@export GOOS=darwin && export GOARCH=amd64 && export CGO_ENABLED=0 && \
		go build -ldflags "-s -w" -o supervisord-monitor

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f supervisord-monitor supervisord-monitor.exe
	@cd frontend && rm -rf node_modules dist

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@cd frontend && npm install

run: ## Run the application
	@go run main.go
