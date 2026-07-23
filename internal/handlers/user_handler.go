package handlers

import (
	"log"
	"net/http"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


type CreateUserHandlerType struct {

	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUserHandler(pool *pgxpool.Pool ) gin.HandlerFunc {
	return func (c *gin.Context){
		
		var input CreateUserHandlerType;

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind JSON"})
			log.Println("bind JSON failed:", err.Error())
			return
		}

		if len(input.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error":  " Password must be more than 6 characters!"})
			return 
		}

		hashedpassword,err := bcrypt.GenerateFromPassword([]byte(input.Password),bcrypt.DefaultCost )

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			log.Println("hash password failed:", err.Error())
			return
		}

		users := &models.User{
			Email: input.Email,
			Password: string(hashedpassword),
		}

		user,err := repository.CreateUser(pool, users)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": " failed to create user"})
			log.Println("create user failed:", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{ "User Created successfully": user})
	}
}