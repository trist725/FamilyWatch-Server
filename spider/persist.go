package spider

type CrawlResult struct {
	Url   string
	Title string
	Img   string
	Dur   string
}

func Persistence(crawResults []CrawlResult) {
	//if _, err := mongo.Collection.UpdateMany(); err != nil {
	//
	//}
}
