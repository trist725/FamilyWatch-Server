package main

import (
	"FamilyWatch/conf"
	"FamilyWatch/db/mongo"
	"FamilyWatch/global"
	myspider "FamilyWatch/spider"
	"FamilyWatch/ws"
	"github.com/hu17889/go_spider/core/spider"
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
			qqSpider.AddUrl(url, "html")
			qqSpider.Run()
		}
	} else {
		//todo: 读库

	}

	go global.Sync()

	ws.Start()
}
