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
	"go.mongodb.org/mongo-driver/mongo/options"
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
func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie id is required"})
			return
		}

		// Expected body
		var req struct {
			AdminReview string `json:"admin_review"`
			RankingName string `json:"ranking_name"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		updateFields := bson.M{}
		if req.AdminReview != "" {
			updateFields["admin_review"] = req.AdminReview
		}
		if req.RankingName != "" {
			updateFields["ranking.ranking_name"] = req.RankingName
		}

		if len(updateFields) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		update := bson.M{"$set": updateFields}

		result := movieCollection.FindOneAndUpdate(
			ctx,
			bson.M{"imdb_id": movieId},
			update,
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		)

		var updatedMovie models.Movie
		if err := result.Decode(&updatedMovie); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Movie updated successfully",
			"movie":   updatedMovie,
		})
	}
}
