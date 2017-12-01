package wx

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	API_GETTICKET_URL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// Js struct
type Js struct {
	context    *WxPublicApplication
	ticket     *resTicket
	updateTime time.Time
}

// Config 返回给用户jssdk配置信息
type Config struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

// resTicket 请求jsapi_tikcet返回结果
type resTicket struct {
	CommonError

	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

//GetConfig 获取jssdk需要的配置参数
//uri 为当前网页地址
func (js *Js) GetConfig(uri string) (config *Config, err error) {
	config = new(Config)
	var ticketStr string
	ticketStr, err = js.GetTicket()
	if err != nil {
		return
	}

	nonceStr := RandomStr(16)
	timestamp := GetCurrTimeStamps()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	sigStr := Signature(str)

	config.AppID = js.context.AppId
	config.NonceStr = nonceStr
	config.Timestamp = timestamp
	config.Signature = sigStr
	return
}

//GetTicket 获取jsapi_tocket
func (js *Js) GetTicket() (ticketStr string, err error) {
	if js.ticket == nil {
		now := time.Now()
		var ticket resTicket
		ticket, err = js.getTicketFromServer()

		if err != nil {
			return
		}
		js.ticket = &ticket
		js.updateTime = now
		ticketStr = ticket.Ticket
	} else {
		nao := time.Duration(7200) * time.Second
		if js.updateTime.Add(nao).Before(time.Now()) {
			js.ticket = nil
			return js.GetTicket()
		}
		ticketStr = js.ticket.Ticket
	}

	return
}

//getTicketFromServer 强制从服务器中获取ticket
func (js *Js) getTicketFromServer() (ticket resTicket, err error) {
	var accessToken string
	accessToken, err = js.context.GetAccessToken()
	if err != nil {
		return
	}

	var response []byte
	url := fmt.Sprintf(API_GETTICKET_URL, accessToken)
	response, err = HTTPGet(url)
	fmt.Println(string(response))
	err = json.Unmarshal(response, &ticket)
	if err != nil {
		return
	}
	if ticket.ErrCode != 0 {
		err = fmt.Errorf("getTicket Error : errcode=%d , errmsg=%s", ticket.ErrCode, ticket.ErrMsg)
		return
	}

	return
}
