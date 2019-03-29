package mongo

import (
	"FamilyWatch/conf"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	CrawledColl *mongo.Collection
	UserColl    *mongo.Collection
)

func Init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Conf.MongoURI))

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping mongo: %v", err)
		return
	}

	CrawledColl = client.Database("FamilyWatch").Collection("crawled")
	UserColl = client.Database("FamilyWatch").Collection("user")
}

func Dispose() {
	CrawledColl = nil
	UserColl = nil
}
