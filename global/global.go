package global

import (
	"FamilyWatch/db/mongo"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CrawlResult struct {
	//_id
	Id       string `bson:"_id"`
	Url      string
	Title    string
	Img      string
	Dur      string
	RealPath string
	Vid      string
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

type QqGetInfo struct {
	Fl fl
	Vl struct {
		Vi []struct {
			Cl struct {
				Fc    string
				Keyid string
			} `json:"cl"`
			Ul struct {
				Ui []struct {
					Url string
				}
			}
			Fn    string
			Fvkey string
		} `json:"vi, omitempty"`
	}
}

//todo: 用于清晰度
type fl struct {
	Fi []struct {
		Cname string
	} `json:"fi"`
}

func GetRealPath(vid string) string {
	if vid == "" {
		log.Print("invalid vid..")
		return ""
	}

	var (
		getInfo   QqGetInfo
		platforms = []string{"1", "4100201", "11", "101100"}
		ret       string
	)
	for _, p := range platforms {
		getInfoUrl := "https://vv.video.qq.com/getinfo?otype=json&platform=" + p + "&defnpayver=1&vid=" + vid
		resp, _ := http.Get(getInfoUrl)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("getInfo failed..")
			return ""
		}
		defer resp.Body.Close()
		body = body[len("QZOutputJson=") : len(body)-1]
		json.Unmarshal(body, &getInfo)

		host := getInfo.Vl.Vi[0].Ul.Ui[0].Url
		fileName := getInfo.Vl.Vi[0].Fn
		vkey := getInfo.Vl.Vi[0].Fvkey

		ret = host + fileName + "?vkey=" + vkey
		return ret
	}

	return ""
}
