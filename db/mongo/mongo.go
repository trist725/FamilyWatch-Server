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

	KaiXinColl   *mongo.Collection
	XiaoDaoColl  *mongo.Collection
	UserColl     *mongo.Collection
	LvYouColl    *mongo.Collection
	JianKangColl *mongo.Collection
	ZhuFuColl    *mongo.Collection
)

func Init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Conf.MongoURI))

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping mongo: %v", err)
		return
	}

	LvYouColl = client.Database("FamilyWatch").Collection("LvYou")
	JianKangColl = client.Database("FamilyWatch").Collection("JianKang")
	KaiXinColl = client.Database("FamilyWatch").Collection("KaiXin")
	XiaoDaoColl = client.Database("FamilyWatch").Collection("XiaoDao")
	ZhuFuColl = client.Database("FamilyWatch").Collection("ZhuFu")
	UserColl = client.Database("FamilyWatch").Collection("user")

	Colls = append(Colls, KaiXinColl)
	Colls = append(Colls, XiaoDaoColl)
	Colls = append(Colls, LvYouColl)
	Colls = append(Colls, JianKangColl)
	Colls = append(Colls, ZhuFuColl)
	Colls = append(Colls, UserColl)

}

func Dispose() {
	for _, c := range Colls {
		c = nil
		_ = c
	}
}
