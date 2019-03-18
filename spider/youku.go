package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
)

type YoukuResult []CrawlResult

type YoukuVideoPageProcessor struct {
}

func NewYoukuVideoPageProcessor() *YoukuVideoPageProcessor {
	return &YoukuVideoPageProcessor{}
}

func (this *YoukuVideoPageProcessor) Process(p *page.Page) {
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

	})

	p.AddTargetRequest(band, "html")

}

func (this *YoukuVideoPageProcessor) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}
