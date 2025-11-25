package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Example route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "MagicStream Movies API is running 🚀",
		})
	})

	// Start server on port 8080
	router.Run(":8080")
}
