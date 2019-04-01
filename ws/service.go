package ws

import (
	"FamilyWatch/conf"
	"FamilyWatch/db/mongo"
	"FamilyWatch/global"
	"context"
	"encoding/json"
	"flag"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func Start() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		serveWss(hub, w, r)
	})
	err := http.ListenAndServe(conf.Conf.WsAddr, nil)
	//err := http.ListenAndServeTLS(conf.Conf.WsAddr,
	//	conf.Conf.CertFile,
	//	conf.Conf.KeyFile, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (req Request) Process() (resp Respond) {
	var (
		restResp *http.Response
		w2sJson  W2S_Code2Session
		userData = global.Users
		userTmp  = global.User{}
		ctx, _   = context.WithTimeout(context.Background(), 10*time.Second)
	)

	if req.Openid != "" {
		if _, ok := userData[req.Openid]; !ok {
			//不在内存中则查找是否在数据库中
			if err := mongo.UserColl.FindOne(ctx, bson.M{"_id": req.Openid}).Decode(&userTmp); err != nil {
				log.Print("UserColl.FindOne: ", err)
				//数据库也找不到
				if req.Code == "" {
					resp.Errcode = 1
					return resp
				} else {
					goto STARTOP
				}
			}
			//数据库里有,读入内存
			userData[req.Openid] = &userTmp
		}
		userTmp.LastLogin = time.Now().Unix()
	}

STARTOP:
	switch req.Op {
	case 1:
		restResp, _ = http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" +
			"wx3cdf9c0b5acf3e86" + "&secret=" + "7cb1a56dbea22e07ca6ef24999abdc97" +
			"&js_code=" + req.Code + "&grant_type=" + "authorization_code")
		defer restResp.Body.Close()
		if err := json.NewDecoder(restResp.Body).Decode(&w2sJson); err != nil {
			log.Print("code2session_decode: ", err)
			resp.Errcode = 1
			break
		}
		if w2sJson.Errcode != 0 {
			resp.Errcode = 1
			log.Print("code2session_errcode: ", w2sJson.Errcode)
			break
		}
		resp.Openid = w2sJson.Openid

	case 2:
		_, ok := userData[req.Openid]
		if !ok {
			log.Print("get user failed on case 2")
			resp.Errcode = 2
			return resp
		}

		switch req.Rcategory {
		case "随机":
		case "养生":
		case "佛缘":
		case "孝道":
		default:
			resp.Errcode = 2
			break
		}

	case 3:
		_, ok := userData[req.Openid]
		if !ok {
			log.Print("get user failed on case 3")
			resp.Errcode = 3
			return resp
		}

	default:

	}

	resp.Op = req.Op

	return resp
}
