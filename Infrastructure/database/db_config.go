package database

import (
	"context"
	"g6/blog-api/Configs"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(env *Configs.Env) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(env.ContextTimeout)*time.Second)
	defer cancel()

	DBHostURI := env.MongoUri
	clientOptions := options.Client().ApplyURI(DBHostURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client *mongo.Client) {
	if client == nil {
		log.Println("MongoDB client is nil, nothing to close.")
		return
	}
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
