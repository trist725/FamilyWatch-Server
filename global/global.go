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
	"strings"
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
	Faved    bool
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
	QQCrawled    = make(map[string]map[string]*CrawlResult)
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
				Fc    int
				Keyid string
				Ci    string
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

type QqGetKey struct {
	Key string
}

func GetRealPath(vid string) string {
	if vid == "" {
		log.Print("invalid vid..")
		return ""
	}

	var (
		getInfo      QqGetInfo
		getKey       QqGetKey
		platforms    = []string{"1", "4100201", "11", "101100"}
		ret          string
		partFormatId string
	)
	for _, p := range platforms {
		getInfoUrl := "https://vv.video.qq.com/getinfo?otype=json&platform=" + p + "&defnpayver=1&vid=" + vid
		respInfo, _ := http.Get(getInfoUrl)
		bodyInfo, err := ioutil.ReadAll(respInfo.Body)
		if err != nil {
			log.Print("getInfo failed..")
			return ""
		}
		defer respInfo.Body.Close()
		if len(bodyInfo) <= (len("QZOutputJson=") + 1) {
			log.Print("invalid bodyInfo")
			return ""
		}
		bodyInfo = bodyInfo[len("QZOutputJson=") : len(bodyInfo)-1]
		json.Unmarshal(bodyInfo, &getInfo)

		if len(getInfo.Vl.Vi) <= 0 {
			log.Print("getInfo vi slice len too short")
			return ""
		} else if len(getInfo.Vl.Vi[0].Ul.Ui) <= 0 {
			log.Print("getInfo ui slice len too short")
			return ""
		}
		host := getInfo.Vl.Vi[0].Ul.Ui[0].Url
		fileName := getInfo.Vl.Vi[0].Fn
		vkey := getInfo.Vl.Vi[0].Fvkey

		ret = host + fileName + "?vkey=" + vkey

		if getInfo.Vl.Vi[0].Cl.Fc == 0 {
			part_format_id := getInfo.Vl.Vi[0].Cl.Keyid
			sp := strings.Split(part_format_id, ".")
			partFormatId = sp[len(sp)-1]
		} else {
			//part_format_id := getInfo.Vl.Vi[0].Cl.Ci[]
		}

		getKeyUrl := "https://vv.video.qq.com/getkey?otype=json&platform=11&format=" + partFormatId + "&vid=" + vid + "&filename=" + fileName
		respKey, _ := http.Get(getKeyUrl)
		bodyKey, err := ioutil.ReadAll(respKey.Body)
		if err != nil {
			log.Print("getKey failed..")
			return ""
		}
		defer respKey.Body.Close()
		if len(bodyKey) <= (len("QZOutputJson=") + 1) {
			log.Print("invalid bodyKey")
			return ret
		}
		bodyKey = bodyKey[len("QZOutputJson=") : len(bodyKey)-1]
		json.Unmarshal(bodyKey, &getKey)
		key := getKey.Key

		ret = host + fileName + "?vkey=" + key

		return ret
	}

	return ""
}
