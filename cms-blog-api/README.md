# CMS Blog Management API (Golang + Gin + MySQL + Redis)

Welcome! Even if you are completely new to Golang, this README will explain every part of the application, how the architecture works, and how to run & test the API step-by-step.

---

## 🏛️ Project Architecture Explained (Beginner Guide)

This application is built using a **Layered Architecture**:

```
HTTP Request
     │
     ▼
┌─────────────────────────┐
│ Middleware (Logger)     │  Logs HTTP Method, Path, Status & Processing Time
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│ Routes (routes.go)      │  Maps URL paths (/blogs) to Controller handlers
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│ Controller              │  Parses JSON / URL params, validates request data
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│ Service Layer           │  Checks Redis Cache (5-min TTL) first!
└────────────┬────────────┘
             │
      ┌──────┴──────┐
      ▼             ▼
┌───────────┐ ┌───────────┐
│   Redis   │ │ Repository│  Queries MySQL via GORM
└───────────┘ └─────┬─────┘
                    ▼
              ┌───────────┐
              │   MySQL   │
              └───────────┘
```

### Key Folders & Responsibilities
- `cmd/main.go` — The entry point of the app. Initializes MySQL, Redis, routes, and starts the HTTP server.
- `config/` — Connects to MySQL (`database.go`) and auto-migrates database tables (`models.Blog`), plus connects to Redis (`redis.go`).
- `models/blog.go` — Defines the `Blog` struct with database columns (`gorm`) and validation tags (`binding:"required"`).
- `controllers/` — Receives HTTP requests, validates input JSON, and sends structured JSON responses.
- `services/` — Business logic layer that checks Redis before hitting MySQL and automatically invalidates cache on updates.
- `repositories/` — Directly interacts with MySQL using GORM queries.
- `middleware/logger.go` — Logs custom performance lines like: `POST /blogs -> 201 Created (20ms)`.

---

## 🚀 How to Run Locally

### Step 1: Start MySQL and Redis using Docker Compose
Open your terminal inside `cms-blog-api/` and run:

```bash
docker compose up -d
```

This starts:
- MySQL 8.0 on port `3306` (Database: `cms_blog_db`, Password: `secret`)
- Redis on port `6379`

### Step 2: Download Go Dependencies
Run:

```bash
go mod tidy
```

### Step 3: Run the API Server
Run:

```bash
go run cmd/main.go
```

You will see:
```text
MySQL Database connected successfully
Database auto-migration completed successfully
Redis connected successfully
Starting CMS Blog API Server on port 8080...
```

---

## 🧪 Testing the APIs with cURL

### 1. Create a new Blog Post (`POST /blogs`)
```bash
curl -X POST http://localhost:8080/blogs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with Golang",
    "content": "Golang is fast, simple, and powerful!",
    "author": "John Doe",
    "category": "Programming",
    "status": "published"
  }'
```

### 2. Get All Published Blogs (`GET /blogs`) - Cached in Redis
```bash
curl http://localhost:8080/blogs
```
*(First call hits MySQL and caches in Redis for 5 minutes. Subsequent calls within 5 minutes return instantly from Redis!)*

### 3. Get Blog by ID (`GET /blogs/1`) - Cached in Redis
```bash
curl http://localhost:8080/blogs/1
```

### 4. Get Blogs by Category (`GET /blogs/category/Programming`) - Cached in Redis
```bash
curl http://localhost:8080/blogs/category/Programming
```

### 5. Update Blog (`PUT /blogs/1`) - Invalidates Cache
```bash
curl -X PUT http://localhost:8080/blogs/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with Golang (Updated)",
    "content": "Updated content goes here...",
    "author": "John Doe",
    "category": "Programming",
    "status": "published"
  }'
```

### 6. Delete Blog (`DELETE /blogs/1`) - Invalidates Cache
```bash
curl -X DELETE http://localhost:8080/blogs/1
```
