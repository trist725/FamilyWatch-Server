package spider

import (
	mymongo "FamilyWatch/db/mongo"
	"FamilyWatch/global"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func Persistence(category string, crawResults []*global.CrawlResult) {
	var (
		coll     *mongo.Collection
		crawlTmp global.CrawlResult
	)

	switch category {
	case "佛缘":
		coll = mymongo.FoyuanColl
	case "孝道":
		coll = mymongo.XiaoDaoColl
	case "养生":
		coll = mymongo.YangShengColl
	default:
		return
	}

	for _, cr := range crawResults {
		doc := bson.M{
			"url":   cr.Url,
			"dur":   cr.Dur,
			"img":   cr.Img,
			"title": cr.Title,
		}
		crawlTmp.Id = ""
		if err := coll.FindOne(context.Background(), bson.D{{"url", cr.Url}}).Decode(&crawlTmp); err == nil {
			//已存在
			if crawlTmp.Id != "" {
				continue
			}
		}

		rId, err := coll.InsertOne(context.Background(), doc)
		if err != nil {
			log.Print("insert one: ", err)
			return
		}

		cr.Id = rId.InsertedID.(primitive.ObjectID).Hex()
	}

	//cur, err := coll.Find(context.Background(), bson.M{})
	//if err != nil {
	//	log.Print("find many: ", err)
	//	return
	//}
	//
	//ctx := context.Background()
	//defer cur.Close(ctx)
	//
	//for cur.Next(ctx) {
	//	elem := global.CrawlResult{}
	//	if err := cur.Decode(&elem); err != nil {
	//		log.Print("cursor decode: ", err)
	//		return
	//	}
	//
	//}
	//
}
