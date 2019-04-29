package main

import (
	"FamilyWatch/conf"
	"FamilyWatch/db/mongo"
	"FamilyWatch/global"
	myspider "FamilyWatch/spider"
	"FamilyWatch/ws"
	"github.com/hu17889/go_spider/core/spider"
	"reflect"
)

func main() {
	mongo.Init()
	defer mongo.Dispose()
	defer global.Save()

	var (
		urls      = conf.Conf.Qq
		runSpider = conf.Conf.RunSpider
		qqSpider  *spider.Spider
	)

	if runSpider {
		for c, url := range urls {
			qqSpider = spider.NewSpider(myspider.NewQQVideoPageProcessor().SetCategory(c), "qqvideo")
			//interface{}转[]string
			s := reflect.ValueOf(url)
			u := make([]interface{}, s.Len())
			for i := 0; i < s.Len(); i++ {
				u[i] = s.Index(i).Interface()
			}

			for i := 0; i < len(u); i++ {
				m := u[i].(string)
				qqSpider.AddUrl(m, "html")
			}

			qqSpider.Run()
		}
	} else {
		//todo: 读库

	}

	go global.Sync()

	ws.Start()
}
