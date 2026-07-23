package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/routes"

	"github.com/gin-gonic/gin"
)

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

	//router allows to routes different HTTP requests to the appropriate handler
	// instead of creating a new instance of gin, i used pointer to get the instance of gin to make my app faster.. it is a Rent vs Buy decision
	// Remember to read the learn.md before continuing...
	// stopped at 2:03:13 ...Database Migration...
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		//gin.H is a map[string]interface{} or map[string]any{}
		ctx.JSON(200, gin.H{
			"message":  "Go running!",
			"status":   "It was successful",
			"database": "connected",
		})
	})

	routes.RegisterTodoRoutes(router, pool)
	routes.RegisterUserRoutes(router, pool)

	router.Run(":" + cfg.Port)
}
