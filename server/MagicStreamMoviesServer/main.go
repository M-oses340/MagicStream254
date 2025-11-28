package main

import (
	"fmt"

	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/database"
	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	router := gin.Default()

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "MagicStream Movies API is running 🚀",
		})
	})
	var client *mongo.Client = database.Connect()
	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	// Start server

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
