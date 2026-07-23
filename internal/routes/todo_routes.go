package routes

import (
	"todo_api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTodoRoutes(router *gin.Engine, pool *pgxpool.Pool) {
	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos", handlers.GetAllTodosHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIdHandler(pool))
	router.PATCH("/todos/:id", handlers.UpdateTodoHandler(pool))
	router.DELETE("/todos/:id", handlers.DeleteTodoHandler(pool))
}
