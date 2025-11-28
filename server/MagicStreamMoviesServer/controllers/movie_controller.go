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

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieCollection := database.OpenCollection("movies")

		cursor, err := movieCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
			return
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

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		movieCollection := database.OpenCollection("movies")
		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
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

		movieCollection := database.OpenCollection("movies")
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()

		role, err := utils.GetRoleFromContext(c)
		if err != nil || role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be part of the ADMIN role"})
			return
		}

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.AdminReview) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin review cannot be empty"})
			return
		}

		// Get sentiment and ranking
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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

		movieCollection := database.OpenCollection("movies")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ranking_name": sentiment,
			"admin_review": req.AdminReview,
		})
	}
}

func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	var sentiments []string
	for _, r := range rankings {
		if r.RankingValue != 999 {
			sentiments = append(sentiments, r.RankingName)
		}
	}
	if len(sentiments) == 0 {
		return "", 0, errors.New("no valid rankings found")
	}

	if err := godotenv.Load(".env"); err != nil {
		// Optional: ignore if env not found
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", 0, errors.New("OPENAI_API_KEY not set")
	}

	basePromptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")
	if basePromptTemplate == "" {
		return "", 0, errors.New("BASE_PROMPT_TEMPLATE not set")
	}

	basePrompt := strings.Replace(basePromptTemplate, "{rankings}", strings.Join(sentiments, ","), 1)

	llm, err := openai.New(openai.WithToken(apiKey))
	if err != nil {
		return "", 0, err
	}

	response, err := llm.Call(context.Background(), basePrompt+adminReview)
	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	for _, r := range rankings {
		if r.RankingName == response {
			rankVal = r.RankingValue
			break
		}
	}

	// Fallback if AI response doesn't match any ranking
	if rankVal == 0 {
		return "Unknown", 0, nil
	}

	return response, rankVal, nil
}
func GetRankings() ([]models.Ranking, error) {
	rankingCollection := database.OpenCollection("rankings")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rankings []models.Ranking
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// panic recovery so we don't bring down the server
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()

		// get user id from context
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		// get user's favourite genres
		favourite_genres, err := GetUsersFavouriteGenres(userId, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// if no favourite genres, return empty list
		if len(favourite_genres) == 0 {
			c.JSON(http.StatusOK, []models.Movie{})
			return
		}

		// load env (optional)
		_ = godotenv.Load(".env")

		// recommended movie limit
		recommendedMovieLimit := int64(5)
		if val := os.Getenv("RECOMMENDED_MOVIE_LIMIT"); val != "" {
			if parsed, perr := strconv.ParseInt(val, 10, 64); perr == nil {
				recommendedMovieLimit = parsed
			}
		}

		// build find options
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(recommendedMovieLimit)

		// build filter
		filter := bson.M{
			"genre": bson.M{
				"$elemMatch": bson.M{
					"genre_name": bson.M{"$in": favourite_genres},
				},
			},
		}

		// query the DB
		movieCollection := database.OpenCollection("movies")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userId string, c *gin.Context) ([]string, error) {
	userCollection := database.OpenCollection("users")

	projection := bson.M{"favourite_genres.genre_name": 1, "_id": 0}
	opts := options.FindOne().SetProjection(projection)

	var result bson.M
	err := userCollection.FindOne(c, bson.D{{Key: "user_id", Value: userId}}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("ℹ️ [GetUsersFavouriteGenres] no user document found for userId=%s\n", userId)
			return []string{}, nil
		}
		log.Printf("❌ [GetUsersFavouriteGenres] FindOne error: %v\n", err)
		return nil, err
	}

	rawGenres, ok := result["favourite_genres"]
	if !ok || rawGenres == nil {
		log.Printf("ℹ️ [GetUsersFavouriteGenres] userId=%s has no favourite genres\n", userId)
		return []string{}, nil
	}

	var genreNames []string

	switch arr := rawGenres.(type) {
	case bson.A:
		for _, item := range arr {
			if genreMap, ok := item.(bson.M); ok {
				if name, ok := genreMap["genre_name"].(string); ok {
					genreNames = append(genreNames, name)
				}
			} else if genreMap, ok := item.(bson.D); ok {
				for _, elem := range genreMap {
					if elem.Key == "genre_name" {
						if name, ok := elem.Value.(string); ok {
							genreNames = append(genreNames, name)
						}
					}
				}
			}
		}
	default:
		log.Printf("⚠️ [GetUsersFavouriteGenres] unexpected type for favourite_genres: %T\n", rawGenres)
		return []string{}, errors.New("unexpected format for favourite_genres")
	}

	log.Printf("ℹ️ [GetUsersFavouriteGenres] userId=%s favourite_genres=%#v\n", userId, genreNames)
	return genreNames, nil
}

func GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		genreCollection := database.OpenCollection("genres")

		cursor, err := genreCollection.Find(c, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		defer cursor.Close(c)

		var genres []models.Genre
		if err := cursor.All(c, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, genres)
	}
}
