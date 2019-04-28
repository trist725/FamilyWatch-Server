package ws

import (
	"FamilyWatch/conf"
	"FamilyWatch/db/mongo"
	"FamilyWatch/global"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func Start() {
	//flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		serveWss(hub, w, r)
	})
	var err error
	if conf.Conf.Encrypt {
		err = http.ListenAndServeTLS(conf.Conf.WsAddr,
			conf.Conf.CertFile,
			conf.Conf.KeyFile, nil)
	} else {
		err = http.ListenAndServe(conf.Conf.WsAddr, nil)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Print("\r\nwebsocket service starting up..")
}

func (req Request) Process() (resp Respond) {
	var (
		restResp    *http.Response
		w2sJson     W2S_Code2Session
		userData    = global.Users
		userTmp     = global.User{}
		ctx, _      = context.WithTimeout(context.Background(), 10*time.Second)
		exist       bool
		reRandCount int
	)

	//刷新直接给数据 不用登陆
	if req.Op == 2 || req.Op == 5 {
		goto STARTOP
	}

	if _, exist = userData[req.Openid]; !exist {
		//不在内存中则查找是否在数据库中
		if err := mongo.UserColl.FindOne(ctx, bson.M{"_id": req.Openid}).Decode(&userTmp); err != nil {
			log.Print("UserColl.FindOne: ", err)
			//数据库也找不到
			if req.Code == "" || req.Op != 1 {
				resp.Errcode = -1
				goto END
			} else {
				goto STARTOP
			}
		}
		//数据库里有,读入内存
		userData[req.Openid] = &userTmp
	}

	userTmp.LastLogin = time.Now().Unix()

STARTOP:
	switch req.Op {
	case 1:
		if exist {
			break
		}
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
		//写入内存待同步
		if _, exist = userData[w2sJson.Openid]; !exist {
			userData[w2sJson.Openid] = &global.User{
				Openid:     w2sJson.Openid,
				SessionKey: w2sJson.Session_key,
				Unionid:    w2sJson.Unionid,
				LastLogin:  time.Now().Unix(),
			}
		}

	case 2:
		if req.Rnum <= 0 || req.Rnum > conf.Conf.RefreshLimit {
			resp.Errcode = 2
			break
		}
		var crawled map[string]*global.CrawlResult
		var ok bool
		if req.Rcategory == "随机" {
			for _, c := range global.QQCrawled {
				crawled = c
				break
			}
		} else {
			if crawled, ok = global.QQCrawled[req.Rcategory]; !ok {
				resp.Errcode = 2
				break
			}
		}

		for ; req.Rnum > 0; req.Rnum-- {
			if len(crawled) == 0 {
				break
			}
			//r := util.RandomInt(0, len(crawled))
			//去重
		RERAND:
			if reRandCount > 100 {
				resp.Errcode = -2
				resp.Load = req.Load
				goto END
			}
			var cTmp *global.CrawlResult
			for _, cTmp = range crawled {
				break
			}
			for _, res := range resp.Resources {
				//这里用url比会出现http和https的相同资源
				if cTmp.Vid == res.Vid {
					reRandCount++
					goto RERAND
				}
			}
			//实时获取真实地址
			if cTmp != nil {
				cTmp.RealPath = global.GetRealPath(cTmp.Vid)
				resp.Resources = append(resp.Resources, *cTmp)
			}
		}

		if f, ok := userData[req.Openid]; ok {
			resp.Favs = f.Favs
		}

		resp.Errcode = 0
		resp.Load = req.Load

	case 3:
		for i, f := range userData[req.Openid].Favs {
			//已收藏
			if f == req.FavId {
				//取消收藏
				userData[req.Openid].Favs = append(userData[req.Openid].Favs[:i], userData[req.Openid].Favs[i+1:]...)
				goto END
			}
		}
		for _, c := range global.QQCrawled {
			for _, v := range c {
				if v.Vid == req.FavId {
					userData[req.Openid].Favs = append(userData[req.Openid].Favs, v.Vid)
					resp.Errcode = 0
					goto END
				}
			}
		}
		resp.Errcode = 3
	case 4:
		resp.Favs = userData[req.Openid].Favs
		for _, f := range userData[req.Openid].Favs {
			for _, c := range global.QQCrawled {
				for vid, v := range c {
					if vid == f {
						resp.Resources = append(resp.Resources, *v)
					}
				}
			}
		}

		resp.Errcode = 0
	case 5:
		if req.Vid == "" {
			resp.Errcode = 5
			return
		}
		resp.RealPath = global.GetRealPath(req.Vid)
		resp.Errcode = 0

	default:
		log.Print("invalid op")
		resp.Errcode = -1
	}

END:
	resp.Op = req.Op

	return resp
}
