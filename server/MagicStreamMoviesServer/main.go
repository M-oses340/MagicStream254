package main

import (
	controller "github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/controllers"
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

	// Movies route
	router.GET("/movies", controller.GetMovies)

	// Start server
	router.Run(":8080")
}
