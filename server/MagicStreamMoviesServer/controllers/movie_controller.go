package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/database"
	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var movieCollection = database.OpenCollection("movies")
var validate = validator.New()

func GetMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var movies []models.Movie

	cursor, err := movieCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
func GetMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movieID := c.Param("imdb_id")
	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "movie id is required"})
		return
	}

	var movie models.Movie
	err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}
func AddMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var movie models.Movie

	// Bind JSON body
	if err := c.BindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Assign new ObjectID
	movie.ID = primitive.NewObjectID()

	// Insert into MongoDB
	_, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movie added successfully", "movie": movie})
}
