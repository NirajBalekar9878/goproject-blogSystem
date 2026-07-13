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
- `config/` — Connects to MySQL (`database.go`) and auto-migrates database tables (`models.Blog`, `models.User`), plus connects to Redis (`redis.go`).
- `models/` — Defines the `Blog` and `User` structs with database columns (`gorm`) and validation tags (`binding:"required"`).
- `controllers/` — Receives HTTP requests (`BlogController`, `AuthController`), validates input JSON, and sends structured JSON responses.
- `services/` — Business logic layer for blogs (Redis caching & MySQL) and authentication (bcrypt password hashing & login).
- `repositories/` — Directly interacts with MySQL using GORM queries (`BlogRepository`, `UserRepository`).
- `middleware/` — Contains custom middleware:
  - `auth_middleware.go` — Validates `Authorization: Bearer <token>` headers and protects private routes.
  - `logger.go` — Logs custom performance lines like: `POST /blogs -> 201 Created (20ms)`.
- `utils/` — Contains JWT token generation & verification (`jwt.go`) and JSON response formatters (`response.go`).

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

### A. Authentication & JWT Endpoints

#### 1. Register a New User (`POST /auth/register`)
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "secretpassword",
    "role": "admin"
  }'
```

#### 2. Login & Receive JWT Token (`POST /auth/login`)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secretpassword"
  }'
```
*Response Example:*
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "role": "admin"
    }
  }
}
```

#### 3. Get Authenticated User Profile (`GET /auth/profile` - Protected)
Set your token in the `Authorization` header:
```bash
TOKEN="YOUR_JWT_TOKEN_HERE"

curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

### B. Blog Management Endpoints

#### 4. Create a new Blog Post (`POST /blogs` - Protected)
```bash
TOKEN="YOUR_JWT_TOKEN_HERE"

curl -X POST http://localhost:8080/blogs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with Golang",
    "content": "Golang is fast, simple, and powerful!",
    "author": "John Doe",
    "category": "Programming",
    "status": "published"
  }'
```

#### 5. Get All Published Blogs (`GET /blogs` - Public, Cached in Redis)
```bash
curl http://localhost:8080/blogs
```
*(First call hits MySQL and caches in Redis for 5 minutes. Subsequent calls within 5 minutes return instantly from Redis!)*

#### 6. Get Blog by ID (`GET /blogs/1` - Public, Cached in Redis)
```bash
curl http://localhost:8080/blogs/1
```

#### 7. Get Blogs by Category (`GET /blogs/category/Programming` - Public, Cached in Redis)
```bash
curl http://localhost:8080/blogs/category/Programming
```

#### 8. Update Blog (`PUT /blogs/1` - Protected, Invalidates Cache)
```bash
TOKEN="YOUR_JWT_TOKEN_HERE"

curl -X PUT http://localhost:8080/blogs/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with Golang (Updated)",
    "content": "Updated content goes here...",
    "author": "John Doe",
    "category": "Programming",
    "status": "published"
  }'
```

#### 9. Delete Blog (`DELETE /blogs/1` - Protected, Invalidates Cache)
```bash
TOKEN="YOUR_JWT_TOKEN_HERE"

curl -X DELETE http://localhost:8080/blogs/1 \
  -H "Authorization: Bearer $TOKEN"
```

