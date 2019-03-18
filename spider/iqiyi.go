package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
)

type IqiyiResult []CrawlResult

type IqiyiVideoPageProcessor struct {
}

func NewIqiyiVideoPageProcessor() *IqiyiVideoPageProcessor {
	return &IqiyiVideoPageProcessor{}
}

func (this *IqiyiVideoPageProcessor) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}
	var band string
	query := p.GetHtmlParser()
	query.Find(".result_item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		url, _ := s.Find("div>a").Attr("href")
		title, _ := s.Find("div>a>img").Attr("alt")
		img, _ := s.Find("div>a>img").Attr("src")
		//title = strings.Replace(title, " ", "", -1)
		fmt.Printf("Review %d: %s - %s - %s\n", i, url, title, img)
		//strings.fi()
	})

	p.AddTargetRequest(band, "html")

}

func (this *IqiyiVideoPageProcessor) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}
