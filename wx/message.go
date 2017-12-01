package wx

import "errors"

//Text 文本消息
type Text struct {
	CommonToken
	Content string `xml:"Content"`
}

//Voice 语音消息
type Voice struct {
	CommonToken

	Voice struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Voice"`
}

//Video 视频消息
type Video struct {
	CommonToken

	Video struct {
		MediaID     string `xml:"MediaId"`
		Title       string `xml:"Title,omitempty"`
		Description string `xml:"Description,omitempty"`
	} `xml:"Video"`
}

//Music 音乐消息
type Music struct {
	CommonToken

	Music struct {
		Title        string `xml:"Title"        `
		Description  string `xml:"Description"  `
		MusicURL     string `xml:"MusicUrl"     `
		HQMusicURL   string `xml:"HQMusicUrl"   `
		ThumbMediaID string `xml:"ThumbMediaId"`
	} `xml:"Music"`
}

//Image 图片消息
type Image struct {
	CommonToken

	Image struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Image"`
}

//News 图文消息
type News struct {
	CommonToken

	ArticleCount int        `xml:"ArticleCount"`
	Articles     []*Article `xml:"Articles>item,omitempty"`
}

//Article 单篇文章
type Article struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicURL      string `xml:"PicUrl,omitempty"`
	URL         string `xml:"Url,omitempty"`
}

//location 地理位置消息
type Location struct {
	CommonToken
	X     float64 `xml:"Location_X"` //地理位置维度
	Y     float64 `xml:"Location_Y"` //地理位置经度
	Scale int     `xml:"Scale"`      //地图缩放大小
	Label string  `xml:"Label"`      //地理位置信息
}

//TransferCustomer 转发客服消息
type TransferCustomer struct {
	CommonToken

	TransInfo *TransInfo `xml:"TransInfo,omitempty"`
}

//TransInfo 转发到指定客服
type TransInfo struct {
	KfAccount string `xml:"KfAccount"`
}

//ErrInvalidReply 无效的回复
var ErrInvalidReply = errors.New("无效的回复消息")

//ErrUnsupportReply 不支持的回复类型
var ErrUnsupportReply = errors.New("不支持的回复消息")

//Reply 消息回复
type Reply struct {
	MsgType MsgType
	MsgData interface{}
}
