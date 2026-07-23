package routes

import (
	"todo_api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(router *gin.Engine, pool *pgxpool.Pool) {
	router.POST("/users", handlers.CreateUserHandler(pool))
}
