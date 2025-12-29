# Bookmark Management API

A scalable and high-performance **Bookmark Management REST API** built with **Golang**.  
This service allows users to create, organize, update, and retrieve bookmarks efficiently, following clean architecture and Go best practices.

---

## ✨ Features

- CRUD operations for bookmarks
- Tag & category support
- Pagination & filtering
- Authentication-ready (JWT-friendly)
- Clean Architecture structure
- RESTful API design
- High-performance and low-memory footprint

---

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```plaintext
.
├── cmd/ # Application entry points
│ └── server/
│ └── main.go
├── internal/
│ ├── handler/
│ ├── model/
│ ├── repository/
│ └── service/
├── pkg/ # Shared utilities
├── configs/ # Configuration files
└── docs/ # API docs, diagrams
```

Benefits:
- Clear separation of concerns
- Easy to test & maintain
- Framework-independent business logic

---

## 🚀 Tech Stack

- **Language:** Go (Golang)
- **HTTP Router:** gin *(configurable)*
- **Database:** PostgreSQL / MySQL *(pluggable)*
- **Auth:** JWT (optional)
- **API Style:** REST
- **Config:** ENV / YAML
- **Docs:** OpenAPI (Swagger-ready)

---

## 🔧 Requirements

- Go **1.21+**
- Docker *(optional but recommended)*
- PostgreSQL / MySQL *(if using database)*

---

## ▶️ Getting Started

### 1️⃣ Clone the repository
```bash
git clone https://github.com/vincent-tien/bookmark-management.git
cd bookmark-management
```

---

## 🛠️ Using Makefile

This project includes a Makefile with convenient commands for common development tasks. All commands should be run from the project root directory.

### Available Commands

#### `make run`
Runs the application directly using `go run`.
```bash
make run
```
This will start the API server on the configured port (default: `8080`).

#### `make swagger`
Generates Swagger/OpenAPI documentation from code annotations.
```bash
make swagger
```
This command uses `swag init` to scan the codebase and generate API documentation files in the `docs/` directory. Run this whenever you update API endpoints or annotations.

#### `make dev-run`
Convenience command that generates Swagger docs and then runs the application.
```bash
make dev-run
```
This is equivalent to running `make swagger` followed by `make run`. Useful for development when you want to ensure docs are up-to-date before starting the server.

#### `make test`
Runs all tests with coverage analysis.
```bash
make test
```
This command:
- Runs tests for all packages (excluding mocks, docs, config, cmd, routers, errors, dto, and test packages)
- Generates a coverage profile (`coverage.out`)
- Creates an HTML coverage report (`coverage.html`)
- Checks if coverage meets the threshold (default: 60%)

**Coverage Options:**
You can customize the coverage threshold and output files using environment variables:
```bash
COVERAGE_THRESHOLD=80 COVERAGE_OUT=mycoverage.out make test
```

**View Coverage Report:**
After running tests, open `coverage.html` in your browser to see a visual coverage report:
```bash
open coverage.html  # macOS
# or
xdg-open coverage.html  # Linux
```

### Example Workflow

```bash
# 1. Generate Swagger docs and start the server
make dev-run

# 2. In another terminal, run tests
make test

# 3. View coverage report
open coverage.html
```

---

## 📝 Additional Setup
