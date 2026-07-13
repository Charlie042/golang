# PostgreSQL Terminal & psql Cheat Sheet
  There are two kinds of "commands" when you use PostgreSQL in the terminal:
  1. **Shell commands** — how you start/connect to Postgres
  2. **psql commands** — what you run *inside* the `psql` client (SQL + backslash meta-commands)
  ---
  ## 1. Starting Postgres from your terminal
  ```bash
  # Connect with defaults (local user, local socket)
  psql

  Connect to a specific database

  psql -d todo_api_g

  Full connection string

  psql "postgres://postgres:password@localhost:5432/todo_api_g?sslmode=disable"

  Connect as a specific user

  psql -U postgres -d todo_api_g -h localhost -p 5432

  Useful flags:
  | Flag | Meaning |
  |------|---------|
  | `-h` | Host |
  | `-p` | Port |
  | `-U` | Username |
  | `-d` | Database name |
  | `-W` | Prompt for password |
  | `-f file.sql` | Run a SQL file |
  | `-c "SQL"` | Run one SQL command and exit |
  ---
  ## 2. Inside `psql`: two types of commands
  ### A. Regular SQL (works everywhere)
  Anything valid in PostgreSQL SQL:
  ```sql
  SELECT * FROM todos;
  INSERT INTO todos (title) VALUES ('Buy milk');
  UPDATE todos SET done = true WHERE id = 1;
  DELETE FROM todos WHERE id = 1;
  CREATE TABLE ...;
  ```
  ### B. psql meta-commands (backslash commands)
  These only work in **psql**, not in Go or most GUI tools.
  Type `\?` inside psql for the full list.
  ---
  ## Connection & session
  | Command | What it does |
  |---------|----------------|
  | `\c dbname` | Connect to another database |
  | `\c dbname user` | Connect as a different user |
  | `\conninfo` | Show current connection info |
  | `\q` | Quit psql |
  | `\password` | Change your password |
  ---
  ## Databases
  | Command | What it does |
  |---------|----------------|
  | `\l` | List all databases |
  | `\l+` | List databases with extra details (size, etc.) |
  ---
  ## Tables, views, schemas
  | Command | What it does |
  |---------|----------------|
  | `\dt` | List tables in current schema |
  | `\dt+` | List tables with size/details |
  | `\dt schema.*` | List tables in a specific schema |
  | `\dv` | List views |
  | `\dm` | List materialized views |
  | `\ds` | List sequences |
  | `\di` | List indexes |
  | `\d table_name` | Describe a table (columns, types, constraints) |
  | `\d+ table_name` | Describe with more detail |
  | `\dn` | List schemas |
  | `\df` | List functions |
  | `\du` | List roles/users |
  | `\du+` | List roles with more detail |
  ---
  ## Query & output formatting
  | Command | What it does |
  |---------|----------------|
  | `\x` | Toggle expanded display (good for wide rows) |
  | `\timing` | Toggle query timing on/off |
  | `\pset` | Set output options (format, borders, etc.) |
  | `\a` | Toggle aligned vs unaligned output |
  | `\H` | HTML output format |
  | `\t` | Show only rows (no headers) |
  | `\g` | Execute current query buffer |
  | `\gx` | Execute with expanded display |
  | `\watch SECS` | Re-run query every N seconds |
  ---
  ## History & editing
  | Command | What it does |
  |---------|----------------|
  | `\s` | Show command history |
  | `\s filename` | Save history to a file |
  | `\e` | Open current query in your editor |
  | `\ef` | Edit a function in your editor |
  ---
  ## Files & scripts
  | Command | What it does |
  |---------|----------------|
  | `\i file.sql` | Run SQL from a file |
  | `\o filename` | Send query output to a file |
  | `\o` | Stop sending output to file (back to terminal) |
  | `\copy` | Import/export data (client-side, different from SQL `COPY`) |
  ---
  ## Help
  | Command | What it does |
  |---------|----------------|
  | `\?` | List all psql meta-commands |
  | `\h` | SQL command help (e.g. `\h SELECT`) |
  | `\h CREATE TABLE` | Help for a specific SQL command |
  ---
  ## Shell integration
  | Command | What it does |
  |---------|----------------|
  | `\! command` | Run a shell command from psql |
  | `\!` | Open a subshell |
  | `\cd path` | Change directory (for `\i`, `\o`, etc.) |
  | `\pwd` | Print current directory |
  ---
  ## Transactions (psql shortcuts)
  | Command | What it does |
  |---------|----------------|
  | `\begin` | Start a transaction |
  | `\commit` | Commit |
  | `\rollback` | Rollback |
  You can also use `BEGIN;`, `COMMIT;`, `ROLLBACK;` as SQL.
  ---
  ## Quick workflow example
  ```bash
  psql -U postgres -d todo_api_g
  ```
  Then inside psql:
  ```text
  \conninfo          -- where am I connected?
  \dt                -- what tables exist?
  \d todos           -- describe the todos table
  SELECT * FROM todos LIMIT 5;
  \q                 -- quit
  ```
  ---
  ## What you can't do in psql meta-commands
  - `\l`, `\dt`, `\d` are **not** SQL — they won't work in your Go app or in pgAdmin's query window the same way.
  - For scripts/automation, use **SQL** or run `psql -c "SELECT ..."` from the shell.
  ---
  ## See everything
  Once you're in psql, run:
  ```text
  \?
  ```
  That prints the complete built-in reference for your installed version.
