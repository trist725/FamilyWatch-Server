package global

type CrawlResult struct {
	//_id
	id    string
	Url   string
	Title string
	Img   string
	Dur   string
	Faved bool
}

type User struct {
	//_id
	Openid     string
	SessionKey string
	Unionid    string
	//存CrawlResult的_id
	Fav       []string
	LastLogin int64
}

var (
	//openid为key
	Users        = make(map[string]*User)
	QQCrawled    = make(map[string][]*CrawlResult)
	IqiyiCrawled = make(map[string][]*CrawlResult)
	YoukuCrawled = make(map[string][]*CrawlResult)
)
