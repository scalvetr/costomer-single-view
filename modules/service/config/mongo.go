package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbConfig struct {
	DbUri      string
	DbUser     string
	DbPassword string
	DbName     string
}

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var MI MongoInstance

func ConnectDB() {

	dbConfig := MongoDbConfig{
		DbUri:      GetEnv("CUSTOMER_DB_URI", "mongodb://localhost:27018"),
		DbName:     GetEnv("CUSTOMER_DB_NAME", "single-customer-view"),
		DbUser:     GetEnv("CUSTOMER_DB_USER", "user"),
		DbPassword: GetEnv("CUSTOMER_DB_PASSWORD", "password"),
	}
	client, err := mongo.NewClient(options.Client().SetAuth(options.Credential{
		AuthSource: dbConfig.DbName,
		Username:   dbConfig.DbUser,
		Password:   dbConfig.DbPassword,
	}).ApplyURI(dbConfig.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected!")

	MI = MongoInstance{
		Client: client,
		DB:     client.Database(dbConfig.DbName),
	}
}
