package handlers

import (
	"net/http"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

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

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := repository.GetAllTodos(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error from failed to get todos": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		errs := repository.DeleteTodo(pool, id)

		if errs != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error from failed to delete todo": errs.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
	}
}
