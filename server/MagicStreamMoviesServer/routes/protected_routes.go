package routes

import (
	controller "github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/controllers"
	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(middleware.AuthMiddleWare())
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
	router.GET("/movie/:imdb_id", controller.GetMovie(client))
	router.POST("/addmovie", controller.AddMovie(client))
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))
}
