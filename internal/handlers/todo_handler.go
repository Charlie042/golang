package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
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

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idStr, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid id format": "Please provide a valid id"})
			return
		}
		todo, err := repository.GetTodoById(pool, idStr)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or the todo does not exist"})
				return
			}
			log.Println("get todo by id failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get todo, there was an internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Your todo was fetched successfully", "status": "success", "data": gin.H{"Todo": todo}})
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		var input UpdateTodoInput

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid id": err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error binding the request body": err.Error()})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please add a title or completed"})
			return
		}
		existing, err := repository.GetTodoById(pool, id)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		title := existing.Title
		if input.Title != nil {
			title = *input.Title
		}

		completed := existing.Completed

		if input.Completed != nil {
			completed = *input.Completed
		}

		todo, errs := repository.UpdateTodo(pool, id, title, completed)
		if errs != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			log.Println("error", errs.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"Your todo has been updated": todo})

	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		errs := repository.DeleteTodo(pool, id)

		if errs != nil {
			if strings.Contains(errs.Error(), "Todo with the id %s not found") {
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error from failed to delete todo": errs.Error()})
			log.Println("error", errs.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
