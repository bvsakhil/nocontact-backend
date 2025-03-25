package main

import (
	"avoids-backend/controllers"
	"avoids-backend/database"
	"avoids-backend/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDatabase()

	// Create Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware (optional)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// User routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	// Avoid routes
	r.POST("/avoids", middleware.AuthMiddleware(), controllers.CreateAvoid)
	r.POST("/check-in", middleware.AuthMiddleware(), controllers.CheckInAvoid)
	r.GET("/avoids", middleware.AuthMiddleware(), controllers.GetUserAvoids)
	r.GET("/avoids/:id", middleware.AuthMiddleware(), controllers.GetAvoidDetails)

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
