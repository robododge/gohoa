package gohoa

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBService struct {
	client     *mongo.Client
	mailboxDB  *mongo.Database
	collection *mongo.Collection
}

func createDBService(collectionName string) DBService {
	config := GetConfig()
	uri := config.MongoDBUrl
	dbName := config.MongoDBName

	uriShort := uri[40:]
	log.Printf(" *--* Connecting to MongoDB: uri [%s]\n", uriShort)
	log.Printf(" *--* Connecting to MongoDB: dbname [%s]\n", dbName)

	ml := DBService{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("  !--! Error connecting to MongoDB: ", err)
	}
	ml.client = client

	//quick ping
	err = ml.client.Ping(ctx, nil)
	if err != nil {
		log.Println(" !--! Error pinging MongoDB: ", err)
	} else {
		log.Println(" *--* Connected to MongoDB")
	}

	ml.mailboxDB = ml.client.Database(dbName)
	ml.collection = ml.mailboxDB.Collection(collectionName)
	return ml

}
