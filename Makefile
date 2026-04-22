# ByteDanceDemo - Parallel Execution Workflow
# Environment setup for China
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off

.PHONY: test test-parallel test-unit test-integration build run-migrate run-api clean

# Run all tests in parallel
test-parallel:
	@echo "Running tests in parallel..."
	@go test -parallel 4 ./... -v 2>&1 | tee test-results.log || true

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -parallel 4 ./... -v -run=Test^[Uu]nit | tee unit-results.log

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -parallel 4 ./... -v -run=Test^[Ii]ntegration | tee integration-results.log

# Run specific test suites
test-service:
	@echo "Running service tests..."
	@go test -parallel 4 ./service -v | tee service-test-results.log

test-repository:
	@echo "Running repository tests..."
	@go test -parallel 4 ./repository -v | tee repository-test-results.log

test-utils:
	@echo "Running utils tests..."
	@go test -parallel 4 ./utils -v | tee utils-test-results.log

# Build the application
build:
	@echo "Building ByteDanceDemo..."
	@go build -o bin/app ./main.go

# Run migration service
run-migrate:
	@echo "Running migration service..."
	@go run cmd/migrate/service.go

# Run API service
run-api:
	@echo "Running API service..."
	@go run cmd/api/service.go

# Run both services in background (parallel)
run-parallel:
	@echo "Starting services in parallel..."
	@go run cmd/migrate/service.go &
	@MIGRATE_PID=$$!
	@go run cmd/api/service.go &
	@API_PID=$$!
	@echo "Migration PID: $$MIGRATE_PID, API PID: $$API_PID"
	@wait

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f bin/app test-results.log unit-results.log integration-results.log
	@rm -f coverage.out coverage.html
	@rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Development workflow - tests and build
dev:
	@echo "Running development workflow..."
	@make test-parallel
	@make build

# CI workflow
ci:
	@echo "Running CI workflow..."
	@make clean
	@make deps
	@make test-parallel
	@make build

# Full test workflow with server
test-full:
	@echo "Running full test workflow..."
	@echo "Starting API server in background..."
	@$(MAKE) run-api > api-server.log 2>&1 & echo $$! > api-server.pid
	@sleep 5
	@echo "Running integration tests..."
	@go test -parallel 4 ./test -v 2>&1 | tee integration-results.log || true
	@echo "Stopping API server..."
	@if [ -f api-server.pid ]; then kill $$(cat api-server.pid) 2>/dev/null || true; rm api-server.pid; fi