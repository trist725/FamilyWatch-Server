package main

import (
	"FamilyWatch/conf"
	"FamilyWatch/db/mongo"
	myspider "FamilyWatch/spider"
	"github.com/hu17889/go_spider/core/spider"
)

func main() {
	mongo.Init()
	defer mongo.Dispose()

	var (
		urls     = conf.Conf.Qq
		qqSpider *spider.Spider
	)

	qqSpider = spider.NewSpider(myspider.NewQQVideoPageProcessor(), "qqvideo").
		//AddPipeline(pipeline.NewPipelineConsole()).
		SetThreadnum(uint(len(urls)))
	for _, url := range urls {
		qqSpider.AddUrl(url, "html")
	}
	qqSpider.Run()

	//spider.NewSpider(spider2.NewQQVideoPageProcessor(), "qqvideo").
	//	AddUrl(req_url, "html"). // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
	//	AddPipeline(pipeline.NewPipelineConsole()).                     // Print result on screen
	//	SetThreadnum(3).                                                // Crawl request by three Coroutines
	//	Run()
	//
	//spider.NewSpider(spider2.NewQQVideoPageProcessor(), "qqvideo").
	//	AddUrl(req_url, "html"). // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
	//	AddPipeline(pipeline.NewPipelineConsole()).                     // Print result on screen
	//	SetThreadnum(3).                                                // Crawl request by three Coroutines
	//	Run()
}
