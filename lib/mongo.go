package lib

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	var err error
	// Connect creates a new Client and then initializes it using the Connect method. This is equivalent to calling NewClient followed by Client.Connect.
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		panic(err)
	}
}

func MongoDBIMDBCollection() *mongo.Collection {
	// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions
	return client.Database("main").Collection("movies")
}

func MongoDBCommentsCollection() *mongo.Collection {
	// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions
	return client.Database("main").Collection("reviews")
}

func MongoDBProcessedReviews() *mongo.Collection {
	// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions
	return client.Database("main").Collection("processed_reviews")
}
