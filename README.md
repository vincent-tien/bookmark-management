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
│ ├── domain/ # Core business logic (entities, interfaces)
│ ├── usecase/ # Application use cases
│ ├── delivery/ # HTTP handlers (REST)
│ └── infrastructure/ # DB, external services, repositories
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
