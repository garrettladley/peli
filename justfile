set dotenv-load := true

# default recipe - show help
default:
    @just --list --unsorted

# ─────────────────────────────────────────────────────────────────────────────
# Development
# ─────────────────────────────────────────────────────────────────────────────

# Install Go dependencies
[group('dev')]
install:
    go mod tidy

# Run the TUI application
[group('dev')]
run:
    go run ./cmd/peli

# Build the application
[group('dev')]
build:
    go build -o peli ./cmd/peli

# Build release binary (stripped)
[group('dev')]
build-release:
    go build -ldflags="-s -w" -o peli ./cmd/peli

# ─────────────────────────────────────────────────────────────────────────────
# Testing
# ─────────────────────────────────────────────────────────────────────────────

# Run all tests
[group('test')]
test:
    go test -v -race ./...

# Run tests with coverage
[group('test')]
test-cover:
    go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# View coverage report in browser
[group('test')]
cover-html: test-cover
    go tool cover -html=coverage.txt

# ─────────────────────────────────────────────────────────────────────────────
# Code Quality
# ─────────────────────────────────────────────────────────────────────────────

# Format code
[group('quality')]
fmt:
    golangci-lint fmt

# Modernize code to use newer Go idioms
[group('quality')]
modernize:
    go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@master -fix ./...

# Run linter
[group('quality')]
lint:
    golangci-lint run --timeout=5m

# Run linter and auto-fix issues
[group('quality')]
lint-fix:
    golangci-lint run --fix --timeout=5m

# Run all code quality checks
[group('quality')]
check: fmt modernize lint test

# ─────────────────────────────────────────────────────────────────────────────
# Database / sqlc
# ─────────────────────────────────────────────────────────────────────────────

migrations_path := "internal/migrations/sqlite/sql"
db_path := env("DB_PATH", "~/.peli/peli.db")

# Generate sqlc code
[group('db')]
sqlc-generate:
    sqlc generate

# Verify sqlc is up to date
[group('db')]
sqlc-diff:
    sqlc diff

# Run sqlc vet
[group('db')]
sqlc-vet:
    sqlc vet

# Create a new migration (up and down files)
[group('db')]
migrate-create name:
    go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -ext sql -dir {{migrations_path}} -seq {{name}}

# Apply all up migrations
[group('db')]
migrate-up:
    go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path={{migrations_path}} -database=sqlite3://{{db_path}} up

# Rollback all migrations
[group('db')]
migrate-down:
    go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path={{migrations_path}} -database=sqlite3://{{db_path}} down

# Rollback last migration
[group('db')]
migrate-down-1:
    go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path={{migrations_path}} -database=sqlite3://{{db_path}} down 1

# Show current migration version
[group('db')]
migrate-version:
    go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path={{migrations_path}} -database=sqlite3://{{db_path}} version

# Force migration version
[group('db')]
migrate-force version:
    go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path={{migrations_path}} -database=sqlite3://{{db_path}} force {{version}}

# Seed database with sample data
[group('db')]
seed:
    go run ./cmd/peli -seed

# ─────────────────────────────────────────────────────────────────────────────
# Maintenance
# ─────────────────────────────────────────────────────────────────────────────

# Clean build artifacts
[group('maint')]
clean:
    rm -f peli coverage.txt
    rm -rf dist/

# Update Go dependencies
[group('maint')]
update:
    go get -u ./...
    go mod tidy
