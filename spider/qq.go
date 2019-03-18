package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
	"strconv"
	"strings"
)

var (
	gCrawlIndex = 1
	gMaxIndex   = 20
	gQQCrawled  QQResult
)

type QQResult []CrawlResult

type QQVideoPageProcessor struct {
}

func NewQQVideoPageProcessor() *QQVideoPageProcessor {
	return &QQVideoPageProcessor{}
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
		//title = strings.Replace(title, " ", "", -1)
		fmt.Printf("Review %d: %s - %s - %s\n", i, url, title, img)
		gQQCrawled = append(gQQCrawled, CrawlResult{
			Url:   url,
			Title: title,
			Img:   img,
		})
	})

	url := p.GetRequest().Url
	old := "&cur=" + strconv.Itoa(gCrawlIndex) + "&"
	gCrawlIndex++
	if gCrawlIndex > gMaxIndex {
		gCrawlIndex = 1
		return
	}
	new := "&cur=" + strconv.Itoa(gCrawlIndex) + "&"
	strings.Replace(url, old, new, 1)
	p.AddTargetRequest(url, "html")
}

func (this *QQVideoPageProcessor) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
	Persistence(gQQCrawled)
}
