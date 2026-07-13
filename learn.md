# What You've Built So Far — Todo API (Go)

A plain-language walkthrough of your codebase. You're at the **foundation stage** of a backend API: the app starts, reads config, connects to Postgres, and serves one test route. The actual "todo" features (create/read/update/delete todos) aren't there yet.

---

## The Big Picture

You're building a **REST API** — a program that listens for HTTP requests (from a browser, Postman, or a frontend) and talks to a **PostgreSQL database**.

Think of it like a restaurant:

| Piece | Your Project | Real-World Analogy |
|-------|--------------|-------------------|
| **Entry point** | `cmd/api/main.go` | Front door — where everything starts |
| **Config** | `internal/config/` | Staff reading the daily setup sheet (.env) |
| **Database** | `internal/database/` | Kitchen — where data is stored |
| **Router (Gin)** | Also in `main.go` | Waiter — takes orders (HTTP requests) and brings responses |
| **`.env`** | Secrets & settings | Private note with DB password, port, etc. |

Right now the waiter only knows one order: `GET /` → "Go running!"

---

## Your Project Structure

```
Golang/
├── cmd/api/main.go          ← App starts here
├── internal/
│   ├── config/config.go     ← Reads .env file
│   └── database/postgres.go ← Connects to Postgres
├── .env                     ← Your secrets (DB URL, port)
├── .air.toml                ← Hot-reload for development
├── go.mod                   ← Go dependencies (like package.json)
└── psql.md                  ← Your Postgres cheat sheet
```

In Go, `cmd/` is the conventional place for **runnable programs**. `internal/` holds code that only this project uses (not meant to be imported by other projects).

---

## What Each File Does

### 1. `cmd/api/main.go` — The Brain

This is where your app boots up. In order, it:

1. **Loads config** from `.env`
2. **Connects to Postgres**
3. **Creates a web server** using Gin (a popular Go HTTP framework)
4. **Defines one route**: `GET /` returns JSON
5. **Starts listening** on port 8080

```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("failed to load config: ", err)
    }

    pool, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("failed to connect to database: ", err)
    }
    defer pool.Close()

    var router *gin.Engine = gin.Default()
    router.SetTrustedProxies(nil)
    router.GET("/", func(ctx *gin.Context) {
        ctx.JSON(200, gin.H{
            "message":  "Go running!",
            "status":   "It was successful",
            "database": "connected",
        })
    })
    router.Run(":" + cfg.Port)
}
```

**Go concepts you're using here:**

- `err` — Go returns errors explicitly instead of throwing exceptions. You check `if err != nil` after almost every operation.
- `defer pool.Close()` — "When this function ends, close the DB connection." Cleanup runs automatically.
- `*gin.Engine` — a pointer. You're reusing Gin's default engine instead of copying it.
- `gin.H{...}` — shorthand for a JSON object (a `map[string]any`).

If you visit `http://localhost:8080/` in a browser, you should see:

```json
{
  "message": "Go running!",
  "status": "It was successful",
  "database": "connected"
}
```

---

### 2. `internal/config/config.go` — Reading Settings

This loads your `.env` file and exposes two values:

```go
type Config struct {
    DatabaseURL string
    Port        string
}

func Load() (*Config, error) {
    var err error = godotenv.Load()

    if err != nil {
        log.Println("Error loading .env file:", err)
        return nil, err
    }
    var config *Config = &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Port:        os.Getenv("PORT"),
    }
    return config, nil
}
```

**Why `.env`?** You don't hardcode passwords in source code. Your `.env` has:

```
DATABASE_URL=postgres://postgres:...@localhost:5432/todo_api_g?sslmode=disable
PORT=8080
```

**Go concepts:**

- `type Config struct` — a custom type grouping related fields (like a TypeScript interface + object shape).
- `*Config` — returns a pointer to the config (common in Go for structs).
- `(*Config, error)` — Go functions often return `(result, error)` as a pair.

---

### 3. `internal/database/postgres.go` — Database Connection

This connects to PostgreSQL using **pgx** (a fast Postgres driver):

```go
func Connect(databaseURL string) (*pgxpool.Pool, error) {
    ctx := context.Background()

    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, err
    }

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, err
    }

    err = pool.Ping(ctx)
    if err != nil {
        pool.Close()
        return nil, err
    }

    log.Println("Database ping successful")
    return pool, nil
}
```

**What's a connection pool?** Instead of opening a new DB connection for every request (slow), the pool keeps several connections ready to reuse.

**Important:** You connect to the DB, but you're **not querying it yet**. No `SELECT`, `INSERT`, or tables for todos. The `/` route just says `"database": "connected"` — it doesn't actually read from Postgres.

---

### 4. `go.mod` — Dependencies

Your main libraries:

| Package | Purpose |
|---------|---------|
| `gin-gonic/gin` | HTTP web framework (routing, JSON responses) |
| `jackc/pgx/v5` | PostgreSQL driver |
| `joho/godotenv` | Load `.env` files |
| `golang-jwt/jwt/v5` | JWT auth — **installed but not used yet** (likely coming in the video) |

---

### 5. `.air.toml` — Developer Convenience

**Air** watches your `.go` files and automatically rebuilds/restarts the server when you save. Similar to `nodemon` in Node.js. Run `air` instead of `go run ./cmd/api` during development.

---

## What Happens When You Start the App

```
You run: air or go run
    ↓
main.go starts
    ↓
config.Load reads .env
    ↓
database.Connect opens Postgres pool
    ↓
Gin router created
    ↓
GET / route registered
    ↓
Server listens on :8080
    ↓
Browser hits localhost:8080
    ↓
Returns JSON response
```

---

## What You've Accomplished

1. Set up a Go module (`todo_api`)
2. Organized code in a standard layout (`cmd/`, `internal/`)
3. Loaded environment variables from `.env`
4. Connected to PostgreSQL with a connection pool
5. Created an HTTP server with Gin
6. Added one health-check style route at `GET /`
7. Set up hot-reload with Air

---

## What's NOT Built Yet

Based on your project name (`todo_api`) and dependencies, the tutorial will likely add:

| Missing Piece | What It Means |
|---------------|---------------|
| **Database tables** | e.g. a `todos` table with `id`, `title`, `done` |
| **Models/structs** | Go types representing a `Todo` |
| **CRUD routes** | `GET /todos`, `POST /todos`, `PUT /todos/:id`, `DELETE /todos/:id` |
| **Handlers** | Functions that run when each route is hit |
| **SQL queries** | Actually reading/writing todos using the `pool` |
| **JWT auth** | Login/register (the jwt package is already in `go.mod`) |
| **Migrations** | Scripts to create/update DB schema |

Right now `pool` is created in `main.go` but never passed to any handler — it's connected and then unused.

---

## Key Backend Ideas to Keep in Mind

### HTTP Methods (you'll use these soon)

- `GET` — read data
- `POST` — create data
- `PUT`/`PATCH` — update data
- `DELETE` — remove data

### The Core Loop

**Request → Handler → Database → Response**

```
Client sends POST /todos { "title": "Buy milk" }
    → Gin routes to your handler
    → Handler runs INSERT INTO todos ...
    → Handler returns JSON { "id": 1, "title": "Buy milk" }
```

### Separation of Concerns

You're already doing this:

- `config` = settings
- `database` = connection logic
- `main` = wiring it together

Later you'll likely add `handlers/`, `models/`, maybe `repository/` for DB queries.

---

## Where You Are in the Tutorial

You're at roughly **"Hello World + database connection"** — the scaffolding before real features. The video probably next covers:

1. Creating the `todos` table in Postgres
2. Defining a `Todo` struct in Go
3. Writing your first real API route (e.g. `GET /todos`)

---

## Quick Reference: Go Syntax You'll See Often

```go
// Variable declaration
var name string = "value"

// Short declaration (inside functions)
name := "value"

// Error handling
result, err := doSomething()
if err != nil {
    log.Fatal(err)
}

// Struct
type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

// Function returning multiple values
func Load() (*Config, error) {
    // ...
    return config, nil  // nil means no error
}

// Defer (run when function exits)
defer pool.Close()
```
