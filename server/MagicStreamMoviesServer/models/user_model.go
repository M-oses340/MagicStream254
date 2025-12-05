package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a MongoDB user document
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"userId"`
	FirstName       string             `bson:"first_name" json:"firstName" validate:"required,min=2,max=100"`
	LastName        string             `bson:"last_name" json:"lastName" validate:"required,min=2,max=100"`
	Email           string             `bson:"email" json:"email" validate:"required,email"`
	Password        string             `bson:"password" json:"password" validate:"required,min=8,max=20"`
	Role            string             `bson:"role" json:"role" validate:"oneof=ADMIN USER"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
	Token           string             `bson:"token" json:"token"`
	RefreshToken    string             `bson:"refresh_token" json:"refreshToken"`
	FavouriteGenres []Genre            `bson:"favourite_genres" json:"favoriteGenres" validate:"required,dive"`
}

// UserLogin represents login input
type UserLogin struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

// UserResponse is used in API responses for user info
type UserResponse struct {
	UserID          string  `json:"userId"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Token           string  `json:"token,omitempty"`
	RefreshToken    string  `json:"refreshToken,omitempty"`
	FavouriteGenres []Genre `json:"favoriteGenres"`
}

// UserContent is the payload inside APIResponse content
type UserContent struct {
	UserID         string  `json:"userId"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	Token          *string `json:"token,omitempty"`
	RefreshToken   *string `json:"refreshToken,omitempty"`
	FavoriteGenres []Genre `json:"favoriteGenres"`
}
