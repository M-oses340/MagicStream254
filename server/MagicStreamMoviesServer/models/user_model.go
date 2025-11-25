package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	FirstName       string             `bson:"first_name" json:"first_name"`
	LastName        string             `bson:"last_name" json:"last_name"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"password"`
	Role            string             `bson:"role" json:"role"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	Token           string             `bson:"token" json:"token"`
	RefreshToken    string             `bson:"refresh_token" json:"refresh_token"`
	FavouriteGenres []Genre            `bson:"favourite_genres" json:"favourite_genres"`
}
