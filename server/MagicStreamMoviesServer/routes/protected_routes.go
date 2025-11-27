package routes

import (
	controller "github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/controllers"
	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())

	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())
}
