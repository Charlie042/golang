# What You've Built So Far — Todo API (Go)

A plain-language walkthrough of your codebase. You've finished the **Todo CRUD stage**: the app starts, reads config, connects to Postgres, and has full create/read/update/delete routes for todos. You're now at the start of the **User Auth stage** — the `users` table migration exists but the register/login/JWT code isn't written yet.

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

The waiter now knows a full menu for todos, but hasn't learned to check IDs at the door yet (no auth) — anyone can order anything.

---

## Your Project Structure

```
Golang/
├── cmd/api/main.go               ← App starts here, routes wired up
├── internal/
│   ├── config/config.go          ← Reads .env file
│   ├── database/postgres.go      ← Connects to Postgres
│   ├── models/todo.go            ← Todo struct (shape of a row)
│   ├── repository/todo_repository.go ← Raw SQL queries for todos
│   └── handlers/todo_handler.go  ← HTTP layer for todos (Create/Get/Update/Delete)
├── migration/                    ← SQL migration files (run via `migrate` CLI)
│   ├── 000001_create_todo_api_table.up/down.sql
│   └── 000002_create_user_api_table.up/down.sql   ← users table, not migrated yet
├── .env                          ← Your secrets (DB URL, port)
├── .air.toml                     ← Hot-reload for development
├── go.mod                        ← Go dependencies (like package.json)
└── psql.md                       ← Your Postgres cheat sheet
```

In Go, `cmd/` is the conventional place for **runnable programs**. `internal/` holds code that only this project uses (not meant to be imported by other projects).

---

## What Each File Does

### 1. `cmd/api/main.go` — The Brain

This is where your app boots up. In order, it:

1. **Loads config** from `.env`
2. **Connects to Postgres**
3. **Creates a web server** using Gin (a popular Go HTTP framework)
4. **Defines a health-check route**: `GET /` returns JSON
5. **Wires up the full todo CRUD routes**, passing the DB `pool` into each handler
6. **Starts listening** on port 8080

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

    router.POST("/todos", handlers.CreateTodoHandler(pool))
    router.GET("/todos", handlers.GetAllTodosHandler(pool))
    router.PATCH("/todos/:id", handlers.UpdateTodoHandler(pool))
    router.GET("/todos/:id", handlers.GetTodoByIdHandler(pool))
    router.DELETE("/todos/:id", handlers.DeleteTodoHandler(pool))
    router.Run(":" + cfg.Port)
}
```

Notice the pattern: `handlers.CreateTodoHandler(pool)` doesn't call the handler directly — it **returns** a `gin.HandlerFunc` that has `pool` baked in via closure. This is how the DB connection gets from `main` down into each handler without a global variable.

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

The `/` route is still just a health check — it doesn't touch the DB. But `pool` is now genuinely used elsewhere: it gets passed into every todo handler.

---

### 4. `internal/models/todo.go` — The Shape of a Todo

```go
type Todo struct {
    ID        int       `json:"id" db:"id"`
    Title     string    `json:"title" db:"title"`
    Completed bool      `json:"completed" db:"completed"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

A **model** is just a Go struct that mirrors a database row. The `json:"..."` tags control how the field is named when it's serialized to JSON in an API response; `db:"..."` documents the matching column name (informational here — pgx maps by position in the `SELECT`/`Scan` calls, not by this tag).

---

### 5. `internal/repository/todo_repository.go` — Talking to Postgres

The **repository layer** is where raw SQL lives. Nothing here knows about HTTP — it just takes a `pool` and plain Go values in, and returns Go structs (or errors) out. Five functions, one per CRUD operation:

- `CreateTodo` — `INSERT ... RETURNING ...`, scans the new row back into a `models.Todo`
- `GetAllTodos` — `SELECT ...`, loops over `rows.Next()` to build a slice
- `GetTodoById` — `SELECT ... WHERE id = $1`, single row
- `UpdateTodo` — `UPDATE ... RETURNING ...`
- `DeleteTodo` — `DELETE ...`, checks `command.RowsAffected()` to know if anything actually matched the id

**Go/pgx concepts:**

- `context.WithTimeout(context.Background(), 5*time.Second)` — every query gets a 5-second budget; if Postgres doesn't respond in time, the query is cancelled instead of hanging forever.
- `$1`, `$2` — placeholders. pgx substitutes these safely (prevents SQL injection) instead of you concatenating strings.
- `.Scan(&todo.ID, &todo.Title, ...)` — copies each returned column into the struct field, in order.

---

### 6. `internal/handlers/todo_handler.go` — Translating HTTP ↔ Repository

The **handler layer** sits between Gin (HTTP) and the repository (SQL). Its job: parse the request, call the repository, shape the response.

```go
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input CreateTodoInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error from binding": err.Error()})
            return
        }
        todo, err := repository.CreateTodo(pool, input.Title, input.Completed)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error from failed to create todo": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, todo)
    }
}
```

Things worth noticing across the five handlers:

- **Input structs** (`CreateTodoInput`, `UpdateTodoInput`) describe exactly what JSON shape is expected. `binding:"required"` makes Gin reject the request automatically if `title` is missing.
- **Pointer fields in `UpdateTodoInput`** (`*string`, `*bool`) let you tell "field not sent" (`nil`) apart from "field sent as empty/false" — needed for a `PATCH` that only updates the fields you actually pass.
- **`errors.Is(err, pgx.ErrNoRows)`** — the standard way to detect "no matching row" and turn it into a `404` instead of a generic `500`.
- Every handler follows the same shape: bind/parse → call repository → map errors to HTTP status codes → respond with JSON.

---

### 7. `migration/` — Database Schema, Versioned

Two migrations exist so far, each with an `.up.sql` (apply) and `.down.sql` (undo):

| File | Creates |
|------|---------|
| `000001_create_todo_api_table` | `todos` table (id, title, description, completed, timestamps) |
| `000002_create_user_api_table` | `users` table (UUID id, email, password, timestamps) — **written but not migrated into your DB yet** |

Migrations are applied with the `golang-migrate` CLI (not a project script — despite what a Windows tutorial reference might suggest with `.\scripts\migrate.ps1`, there's no such script here):

```bash
migrate -database "$DATABASE_URL" -path migration up
```

---

### 8. `go.mod` — Dependencies

Your main libraries:

| Package | Purpose |
|---------|---------|
| `gin-gonic/gin` | HTTP web framework (routing, JSON responses) |
| `jackc/pgx/v5` | PostgreSQL driver |
| `joho/godotenv` | Load `.env` files |
| `golang-jwt/jwt/v5` | JWT auth — **installed but not used yet** (likely coming in the video) |

---

### 9. `.air.toml` — Developer Convenience

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
Gin router created, todo routes registered with pool wired in
    ↓
Server listens on :8080
    ↓
Client hits POST /todos { "title": "Buy milk" }
    ↓
Handler binds JSON → calls repository.CreateTodo → INSERT INTO todos
    ↓
Returns the created todo as JSON
```

---

## What You've Accomplished

1. Set up a Go module (`todo_api`)
2. Organized code in a standard layout (`cmd/`, `internal/`)
3. Loaded environment variables from `.env`
4. Connected to PostgreSQL with a connection pool
5. Created an HTTP server with Gin
6. Set up hot-reload with Air
7. Built a `models` → `repository` → `handlers` layering for todos
8. Implemented full CRUD for todos: `POST /todos`, `GET /todos`, `GET /todos/:id`, `PATCH /todos/:id`, `DELETE /todos/:id`
9. Written (but not yet run) a migration for a `users` table with UUID ids

---

## What's NOT Built Yet

Based on the `README.md` architecture diagram and the `golang-jwt/jwt/v5` dependency already in `go.mod`, the tutorial will likely add next:

| Missing Piece | What It Means |
|---------------|---------------|
| **Run the users migration** | Apply `migration/000002_create_user_api_table.up.sql` — see [psql.md](psql.md) / use `migrate -database "$DATABASE_URL" -path migration up` |
| **User model/repository** | Go struct + SQL for `users`, mirroring how `todo.go`/`todo_repository.go` work |
| **Password hashing** | bcrypt (or similar) before storing `password` — never store plaintext |
| **`POST /auth/register`, `POST /auth/login`** | Public routes to create a user and issue a JWT |
| **Auth middleware** | Validates the JWT on protected routes, attaches the user to the request context |
| **Protecting the todo routes** | Per the README diagram, `/todos/*` should require auth |
| **`user_id` on todos** | The README's DB diagram shows a foreign key from `todos.user_id` → `users.id`; the current `todos` table migration doesn't have this column yet |

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

You've finished the **Todo CRUD milestone** and are at the start of the **User Auth milestone**. Your own note in `main.go` marks this: *"stopped at 2:03:13 ...Database Migration..."*. Concretely:

- ✅ Todos table, model, repository, handlers, routes — all done
- ✅ `users` table migration file written
- ⬜ **Next up:** run the users migration, then build the `users` model/repository, registration/login handlers, JWT issuing, and auth middleware

The immediate next hands-on step is applying `migration/000002_create_user_api_table.up.sql` to your database with the `migrate` CLI (see the migrations section above), then confirming it with:

```bash
migrate -database "$DATABASE_URL" -path migration version
```

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
