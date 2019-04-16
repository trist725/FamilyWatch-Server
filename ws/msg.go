package ws

import "FamilyWatch/global"

type Request struct {
	//操作类型,1-登陆,2-刷新,3-(取消)收藏,4-获取收藏
	Op int `json:"op, omitempty"`
	//客户端调用wx.login()获取到的登陆凭证
	Code string `json:"code, omitempty"`
	//用户唯一标识,必填
	Openid string `json:"openid, omitempty"`
	//刷新类别,"随机","养生"...
	Rcategory string `json:"category, omitempty"`
	//刷新数量
	Rnum int `json:"num, omitempty"`
	//要收藏的视频id
	FavId string `json:"favid, omitempty"`
	//刷新加载,true-刷新,false-加载
	Load bool `json:"load, omitempty"`
}

type Respond struct {
	//操作类型,1-登陆,2-刷新,3-(取消)收藏,4-获取收藏
	Op int `json:"op"`
	//用户唯一标识
	Openid string `json:"openid"`
	//错误码,0-正常,1-openid不对,需要重新登陆,2-刷新参数不对,3-收藏参数不对
	Errcode int `json:"errcode"`
	//资源
	Resources []global.CrawlResult `json:"resources"`
	//已收藏的资源的_id
	Favs []string
	//刷新加载,true-刷新,false-加载
	Load bool `json:"load, omitempty"`
}

type S2W_Code2Session struct {
	appid      string //小程序 appId
	secret     string //小程序 appSecret
	js_code    string //登录时获取的 code
	grant_type string //授权类型，此处只需填写 authorization_code
}

type W2S_Code2Session struct {
	Openid      string //用户唯一标识
	Session_key string //会话密钥
	Unionid     string //用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	Errcode     int    //错误码,0-请求成功,-1-系统繁忙,此时请开发者稍候再试,40029-code无效,45011-频率限制,每个用户每分钟100次
	Errmsg      string //错误信息
}
