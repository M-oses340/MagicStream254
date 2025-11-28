package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set")
	}

	fmt.Println("MONGODB_URI:", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	// Create client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Mongo NewClient Error: %v", err)
	}

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Mongo Connect Error: %v", err)
	}

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Mongo Ping Error: %v", err)
	}

	fmt.Println("MongoDB Connected Successfully!")
	return client
}

var Client *mongo.Client = Connect()

func OpenCollection(collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatal("DATABASE_NAME not set")
	}

	collection := Client.Database(databaseName).Collection(collectionName)
	return collection
}
