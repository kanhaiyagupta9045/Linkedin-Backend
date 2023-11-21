package databases

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() (client *mongo.Client) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while Loading .env files")
	}
	connectionString := os.Getenv("MONGODB_URL")
	if err != nil {
		log.Fatal(err)
	}
	if connectionString == "" {
		log.Fatal("mongo DB error is not present in .env file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	}
	return cli
}

var Client *mongo.Client = DBInstance()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("SocialMedia").Collection(collectionName)
	return collection
}
