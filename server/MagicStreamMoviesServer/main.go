package main

import (
	"fmt"

	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "MagicStream Movies API is running 🚀",
		})
	})
	routes.SetupUnProtectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	// Start server

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
