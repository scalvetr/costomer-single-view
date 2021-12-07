package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
)

const collectionName = "cases"

type MongoDbConfig struct {
	DbUri      string
	DbUser     string
	DbPassword string
	DbName     string
}
type ContactCenterRepo struct {
	dbConfig MongoDbConfig
	client   mongo.Client
}

func BuildContactCenterRepo(dbConfig MongoDbConfig) ContactCenterRepo {
	localClient, err := mongo.NewClient(options.Client().SetAuth(options.Credential{
		Username: dbConfig.DbUser,
		Password: dbConfig.DbPassword,
	}).ApplyURI(dbConfig.DbUri))

	if err != nil {
		panic(err)
	}
	return ContactCenterRepo{
		client:   *localClient,
		dbConfig: dbConfig,
	}
}

func (r ContactCenterRepo) init() (context.Context, mongo.Database, mongo.Collection) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := r.client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer r.client.Disconnect(ctx)
	db := r.client.Database(r.dbConfig.DbName)
	err = db.CreateCollection(ctx, collectionName)
	if err != nil {
		panic(err)
	}
	collection := db.Collection(collectionName)

	return ctx, *db, *collection
}
func (r ContactCenterRepo) GetOpenCase(customerId string) *CaseStruct {
	ctx, _, collection := r.init()

	resultBson, err := collection.Find(ctx, bson.M{"customer_id": customerId})
	if err != nil {
		panic(err)
	}
	if resultBson == nil {
		return nil
	}
	var result []CaseStruct
	resultBson.All(ctx, result)

	if result != nil && len(result) > 0 {
		return &result[rand.Intn(len(result)-1)]
	}
	// no cases for this customer
	return nil
}

func (r ContactCenterRepo) StoreCase(caseStruct CaseStruct) CaseStruct {
	ctx, _, collection := r.init()

	res, err := collection.InsertOne(ctx, caseStruct)
	if err != nil {
		panic(err)
	}

	var result CaseStruct
	err = collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(result)
	if err != nil {
		panic(err)
	}
	return result

}
