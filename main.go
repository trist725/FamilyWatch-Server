package main

import (
	"FamilyWatch/conf"
	"FamilyWatch/db/mongo"
	myspider "FamilyWatch/spider"
	"FamilyWatch/ws"
	"github.com/hu17889/go_spider/core/spider"
)

func main() {
	mongo.Init()
	defer mongo.Dispose()

	var (
		urls      = conf.Conf.Qq
		runSpider = conf.Conf.RunSpider
		qqSpider  *spider.Spider
	)

	if runSpider {
		qqSpider = spider.NewSpider(myspider.NewQQVideoPageProcessor(), "qqvideo").
			//AddPipeline(pipeline.NewPipelineConsole()).
			SetThreadnum(uint(len(urls)))
		for _, url := range urls {
			qqSpider.AddUrl(url, "html")
		}
		go qqSpider.Run()
	}

	ws.Start()
}
