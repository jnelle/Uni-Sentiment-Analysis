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
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		panic(err)
	}
}

func MongoDBGetIMDBCollection() *mongo.Collection {
	return client.Database("uni").Collection("imdb")
}

func MongoDBGetCommentsCollection() *mongo.Collection {
	return client.Database("uni").Collection("comments")
}
