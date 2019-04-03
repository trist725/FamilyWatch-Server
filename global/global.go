package global

type CrawlResult struct {
	//_id
	Id    string `bson:"_id"`
	Url   string
	Title string
	Img   string
	Dur   string
}

type User struct {
	//_id
	Openid     string `bson:"_id"`
	SessionKey string
	Unionid    string
	//存CrawlResult的_id
	Favs      []string
	LastLogin int64
}

var (
	//openid为key
	Users        = make(map[string]*User)
	QQCrawled    = make(map[string][]*CrawlResult)
	IqiyiCrawled = make(map[string][]*CrawlResult)
	YoukuCrawled = make(map[string][]*CrawlResult)
)
