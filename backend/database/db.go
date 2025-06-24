package database

import (
	"context"
	"log"
	"time"

	"f/cloudfrog/backend/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	DB     *mongo.Database
	Files  *mongo.Collection
)

// FileRecord represents a file in the database
type FileRecord struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileName     string             `bson:"fileName" json:"fileName"`
	OriginalName string             `bson:"originalName" json:"originalName"`
	ShortCode    string             `bson:"shortCode" json:"shortCode"`
	MimeType     string             `bson:"mimeType" json:"mimeType"`
	Size         int64              `bson:"size" json:"size"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	ExpiresAt    time.Time          `bson:"expiresAt" json:"expiresAt"`
}

// Connect establishes connection to the MongoDB database
func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(config.MongoDBConnectionString)
	var err error

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	DB = client.Database("cloudfrog")
	Files = DB.Collection("files")

	// Create unique index for shortCode
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"shortCode", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = Files.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	log.Println("MongoDB connected successfully")
}

// Disconnect closes the MongoDB connection
func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}
