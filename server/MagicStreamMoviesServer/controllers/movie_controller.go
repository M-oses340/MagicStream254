package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"movies": []string{"Movie 1", "Movie 2"},
	})
}
