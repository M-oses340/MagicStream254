package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/M-oses340/MagicStream254/server/MagicStreamMoviesServer/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// =========================
// JWT CLAIMS STRUCTURE
// =========================
type SignedDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	UserId    string `json:"user_id"`
	jwt.RegisteredClaims
}

// Load secrets safely
var SECRET_KEY = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")

func init() {
	if SECRET_KEY == "" || SECRET_REFRESH_KEY == "" {
		fmt.Println("⚠ WARNING: JWT SECRET KEYS ARE NOT SET IN ENVIRONMENT VARIABLES")
	}
}

// =========================
// GENERATE ACCESS + REFRESH TOKENS
// =========================
func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {

	// ACCESS TOKEN
	accessClaims := SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 day
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	// REFRESH TOKEN
	refreshClaims := SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

// =========================
// UPDATE TOKENS IN MONGODB
// =========================
func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) error {
	userCollection := database.OpenCollection("users", client)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"token":         token,
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	}}

	_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, update)
	return err
}

// =========================
// EXTRACT 'Authorization: Bearer token'
// =========================
func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return "", errors.New("empty bearer token")
	}

	return tokenString, nil
}

// =========================
// VALIDATE ACCESS TOKEN
// =========================
func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("access token expired")
	}

	return claims, nil
}

// =========================
// GET ROLE FROM TOKEN
// =========================
func GetRoleFromContext(c *gin.Context) (string, error) {
	tokenString, err := GetAccessToken(c)
	if err != nil {
		return "", err
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Role, nil
}

// =========================
// GET USER ID FROM TOKEN
// =========================
func GetUserIdFromContext(c *gin.Context) (string, error) {
	tokenString, err := GetAccessToken(c)
	if err != nil {
		return "", err
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserId, nil
}

// =========================
// VALIDATE REFRESH TOKEN
// =========================
func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_REFRESH_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return claims, nil
}
