package routes

import (
	controller "github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	
	router.GET("/movie/:imdb_id", controller.GetMovie)
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}
