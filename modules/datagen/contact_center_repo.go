package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

func BuildContactCenterRepo(dbConfig MongoDbConfig) ContactCenterRepo {
	return ContactCenterRepo{
		dbConfig: dbConfig,
	}
}

func (r ContactCenterRepo) init() (context.Context, mongo.Database, mongo.Collection) {
	client, err := mongo.NewClient(options.Client().SetAuth(options.Credential{
		AuthSource: r.dbConfig.DbName,
		Username:   r.dbConfig.DbUser,
		Password:   r.dbConfig.DbPassword,
	}).ApplyURI(r.dbConfig.DbUri))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	db := client.Database(r.dbConfig.DbName)

	collectionList, err := db.ListCollections(ctx, bson.M{"name": collectionName})
	if err != nil {
		panic(err)
	}
	if !collectionList.Next(ctx) {
		err = db.CreateCollection(ctx, collectionName)
		if err != nil {
			panic(err)
		}
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

func (r ContactCenterRepo) StoreCase(caseStruct CaseStruct) primitive.ObjectID {
	ctx, _, collection := r.init()

	res, err := collection.InsertOne(ctx, caseStruct)
	if err != nil {
		panic(err)
	}
	return res.InsertedID.(primitive.ObjectID)

}

func (r ContactCenterRepo) Close() error {
	//err := r.client.Disconnect(context.TODO())
	//if err == nil {
	//	fmt.Println("Connection to MongoDB closed.")
	//}
	//return err
	return nil
}
