package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/database"
	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/models"
	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection = database.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return
		}

		c.JSON(http.StatusOK, movies)

	}
}

func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movieCollection = database.OpenCollection("movies", client)

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)

	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		var movieCollection = database.OpenCollection("movies", client)

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}
func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			log.Println("Role not found in context:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}

		log.Println("User role:", role)

		if role != "ADMIN" {
			log.Println("Unauthorized role:", role)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be part of the ADMIN role"})
			return
		}

		// Get movie ID from URL
		movieId := c.Param("imdb_id")
		if movieId == "" {
			log.Println("Movie ID not provided")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id required"})
			return
		}

		// Bind request body
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		log.Println("Received admin review for movie:", movieId, "Review:", req.AdminReview)

		// Get sentiment & ranking
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)
		if err != nil {
			log.Println("GetReviewRanking error:", err, "Input:", req.AdminReview)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}

		log.Println("Sentiment:", sentiment, "Ranking value:", rankVal)

		// Prepare MongoDB update
		filter := bson.D{{Key: "imdb_id", Value: movieId}}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieCollection := database.OpenCollection("movies", client)
		if movieCollection == nil {
			log.Println("Movie collection is nil")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie collection not found"})
			return
		}

		log.Println("Updating movie in MongoDB:", movieId, "with update:", update)

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println("MongoDB UpdateOne error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			log.Println("No movie matched for update with ID:", movieId)
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		log.Println("Movie updated successfully. MatchedCount:", result.MatchedCount, "ModifiedCount:", result.ModifiedCount)

		// Respond with updated data
		c.JSON(http.StatusOK, gin.H{
			"ranking_name": sentiment,
			"admin_review": req.AdminReview,
		})
	}
}

func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	rankings, err := GetRankings(client, c)

	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(c, base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}
	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}
	return response, rankVal, nil

}

func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var rankingCollection = database.OpenCollection("rankings", client)

	cursor, err := rankingCollection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil

}

func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		favourite_genres, err := GetUsersFavouriteGenres(userId, client, c)
		if err != nil {
			log.Println("Error fetching user's favourite genres:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		recommendedMovieLimitVal := int64(5)
		if recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT"); recommendedMovieLimitStr != "" {
			if val, err := strconv.ParseInt(recommendedMovieLimitStr, 10, 64); err == nil {
				recommendedMovieLimitVal = val
			} else {
				log.Println("Error parsing RECOMMENDED_MOVIE_LIMIT:", err)
			}
		}

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(recommendedMovieLimitVal)

		filter := bson.D{
			{Key: "genre.genre_name", Value: bson.D{
				{Key: "$in", Value: favourite_genres},
			}},
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieCollection := database.OpenCollection("movies", client)
		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			log.Println("MongoDB Find error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			log.Println("Error decoding movies cursor:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: userId}}
	projection := bson.M{"favourite_genres.genre_name": 1, "_id": 0}

	var result struct {
		FavouriteGenres []struct {
			GenreName string `bson:"genre_name"`
		} `bson:"favourite_genres"`
	}

	userCollection := database.OpenCollection("users", client)
	err := userCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
		return nil, err
	}

	genreNames := make([]string, 0, len(result.FavouriteGenres))
	for _, g := range result.FavouriteGenres {
		genreNames = append(genreNames, g.GenreName)
	}

	return genreNames, nil
}
func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		genreCollection := database.OpenCollection("genres", client)

		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"error":   true,
				"message": "Error fetching movie genres",
				"content": gin.H{},
			})
			return
		}
		defer cursor.Close(ctx)

		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"error":   true,
				"message": err.Error(),
				"content": gin.H{},
			})
			return
		}

		// Respond with content array
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"error":   false,
			"message": "Genres fetched successfully",
			"content": genres,
		})
	}
}
