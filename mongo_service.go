package gohoa

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbService struct {
	client     *mongo.Client
	mailboxDB  *mongo.Database
	collection *mongo.Collection
}

func createDbService(collectionName string) DbService {
	config := GetConfig()
	uri := config.MongoDbUrl
	dbName := config.MongoDbName

	ml := DbService{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Error connecting to MongoDB: ", err)
	}
	ml.client = client

	//quick ping
	err = ml.client.Ping(ctx, nil)
	if err != nil {
		log.Println("Error pinging MongoDB: ", err)
	} else {
		log.Println("Connected to MongoDB")
	}

	ml.mailboxDB = ml.client.Database(dbName)
	ml.collection = ml.mailboxDB.Collection(collectionName)
	return ml

}
