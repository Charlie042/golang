package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// study about database connection context
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
	INSERT INTO todos (title, completed)
	VALUES ($1, $2)
	RETURNING id, title, completed, created_at, updated_at;
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todos []models.Todo
	rows, err := pool.Query(ctx, `SELECT id, title, completed, created_at, updated_at FROM todos ORDER BY created_at DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoById(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todo models.Todo
	var query = `SELECT id, title, completed, created_at, updated_at FROM todos WHERE id = $1
	`
	if err := pool.QueryRow(ctx, query, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}
	return &todo, nil

}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = ` UPDATE todos SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 
	RETURNING id, title, completed, created_at, updated_at;
	`
	var todo models.Todo
	if err := pool.QueryRow(ctx, query, title, completed, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
	DELETE FROM todos WHERE id = $1;
	`

	command, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return fmt.Errorf("Todo with the id %s not found", id)
	}

	return nil
}
