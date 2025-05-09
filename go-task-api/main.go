package main

import (
	"go-task-api/controllers"
	"go-task-api/middleware"
	"go-task-api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
	}))

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.GET("/me", func(c *gin.Context) {
		user, _ := c.Get("user")
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	protected.POST("/tasks", controllers.CreateTask)
	protected.GET("/tasks", controllers.ListTasks)
	protected.GET("/tasks/:id", controllers.GetTask)
  protected.PUT("/tasks/:id", controllers.UpdateTask)
  protected.DELETE("/tasks/:id", controllers.DeleteTask)
	protected.GET("/trash", controllers.ViewTrashedTasks)
	protected.PUT("/tasks/:id/restore", controllers.RestoreTask)
	protected.DELETE("/tasks/:id/permanent", controllers.HardDeleteTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
