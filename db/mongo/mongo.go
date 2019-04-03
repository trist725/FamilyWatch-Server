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
	Colls []*mongo.Collection

	FoyuanColl    *mongo.Collection
	XiaoDaoColl   *mongo.Collection
	YangShengColl *mongo.Collection
	UserColl      *mongo.Collection
)

func Init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Conf.MongoURI))

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping mongo: %v", err)
		return
	}

	FoyuanColl = client.Database("FamilyWatch").Collection("FoYuan")
	XiaoDaoColl = client.Database("FamilyWatch").Collection("XiaoDao")
	YangShengColl = client.Database("FamilyWatch").Collection("YangSheng")
	UserColl = client.Database("FamilyWatch").Collection("user")

	Colls = append(Colls, FoyuanColl)
	Colls = append(Colls, XiaoDaoColl)
	Colls = append(Colls, YangShengColl)
	Colls = append(Colls, UserColl)

}

func Dispose() {
	for _, c := range Colls {
		c = nil
		_ = c
	}
}
