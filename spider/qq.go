package spider

import (
	"FamilyWatch/conf"
	"FamilyWatch/global"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
	"strconv"
	"strings"
	"time"
)

var (
	gCrawlIndex = 1
)

type QQVideoPageProcessor struct {
	category string
}

func NewQQVideoPageProcessor() *QQVideoPageProcessor {
	return &QQVideoPageProcessor{}
}

func (this *QQVideoPageProcessor) SetCategory(c string) *QQVideoPageProcessor {
	this.category = c
	return this
}

func (this *QQVideoPageProcessor) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()
	query.Find(".result_item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		url, _ := s.Find("div>a").Attr("href")
		title, _ := s.Find("div>a>img").Attr("alt")
		img, _ := s.Find("div>a>img").Attr("src")
		dur := s.Find("div>a>span[class=figure_caption]>span").Text()
		title = strings.Replace(title, "\006", "", -1)
		fmt.Printf("Review %d: %s - %s - %s - %s\n", i, url, title, img, dur)
		min, _ := time.Parse("15:04:05", dur)
		//按分钟过滤
		if min.Minute() >= conf.Conf.FilterMin {
			//关键词分类
			if this.category != "" && strings.Contains(url, "/x/page") {
				global.QQCrawled[this.category] = append(global.QQCrawled[this.category], &global.CrawlResult{
					Url:   url,
					Title: title,
					Img:   img,
					Dur:   dur,
				})
			}
		}
	})

	url := p.GetRequest().Url
	old := "&cur=" + strconv.Itoa(gCrawlIndex) + "&"
	gCrawlIndex++
	if gCrawlIndex > conf.Conf.MaxCrawlIndex {
		gCrawlIndex = 1
		return
	}
	new := "&cur=" + strconv.Itoa(gCrawlIndex) + "&"
	strings.Replace(url, old, new, 1)
	p.AddTargetRequest(url, "html")
}

func (this *QQVideoPageProcessor) Finish() {
	fmt.Printf("[%s] crawled num: [%d] \r\n", this.category, len(global.QQCrawled[this.category]))
	for k, v := range global.QQCrawled {
		Persistence(k, v)
	}
}
