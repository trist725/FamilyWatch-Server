package global

import (
	"FamilyWatch/db/mongo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type CrawlResult struct {
	//_id
	Id    string `bson:"_id"`
	Url   string
	Title string
	Img   string
	Dur   string
}

type User struct {
	//_id
	Openid     string `bson:"_id"`
	SessionKey string
	Unionid    string
	//存CrawlResult的_id
	Favs      []string
	LastLogin int64
}

var (
	//openid为key
	Users        = make(map[string]*User)
	QQCrawled    = make(map[string][]*CrawlResult)
	IqiyiCrawled = make(map[string][]*CrawlResult)
	YoukuCrawled = make(map[string][]*CrawlResult)
)

func Sync() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	for {
		<-t.C
		fmt.Println("Sync...")
		Save()
	}

}

func Save() {
	for _, u := range Users {
		_, err := mongo.UserColl.UpdateOne(context.Background(), bson.M{"_id": u.Openid}, bson.D{{"$set", bson.D{
			{"_id", u.Openid},
			{"SessionKey", u.SessionKey},
			{"Unionid", u.Unionid},
			{"Favs", u.Favs},
			{"LastLogin", u.LastLogin},
		}}}, options.Update().SetUpsert(true))
		if err != nil {
			log.Print("sync: ", err)
		}
	}
}
