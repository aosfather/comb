package wx

import (
	"encoding/xml"
	"strconv"
)

/*
   对消息进行回复响应处理

*/
type ReplyType int

const (
	MSG_Text  ReplyType = 1
	MSG_Image ReplyType = 2
	MSG_Voice ReplyType = 3
	MSG_Video ReplyType = 4
	MSG_Music ReplyType = 5
	MSG_News  ReplyType = 6
)

//微信校验消息
type WxValidateRequest struct {
	Timestamp string `Field:"timestamp"`
	Nonce     string `Field:"nonce"`
	Signature string `Field:"msg_signature"`
	Echostr   string `Field:"echostr"`
}

type ReplyMessage struct {
	Type        ReplyType
	Title       string
	Content     string
	ExtContent  string
	MediaId     string
	Description string
	Items       []ReplyMessageItem
}

type ReplyMessageItem struct {
	Title       string
	Pic         string
	Url         string
	Description string
}

//消息处理
type MessageHandle interface {
	OnTextMessage(user string, text string) ReplyMessage
	OnImageMessage(user string, pic string, mediaId string) ReplyMessage
	OnVoiceMessage(user string, format string, mediaId string, recognition string) ReplyMessage
	OnVideoMessage(user string, thumb string, mediaId string) ReplyMessage
	OnShortVideoMessage(user string, thumb string, mediaId string) ReplyMessage
	OnLocationMessage(user string, label string, x, y float64, scale int) ReplyMessage
	OnLinkMessage(user string, title string, url string, description string) ReplyMessage
}

//事件处理
type EventHandle interface {
	OnSubscribe(user string, unsub bool) ReplyMessage                 //处理订阅事件，unsub为true时候为取消订阅
	OnScan(user string, key int64, ticket string) ReplyMessage        //用户扫描二维码
	OnMenuClick(user string, key string) ReplyMessage                 //菜单项对应的是点击
	OnMenuView(user string, url string) ReplyMessage                  //菜单项对应的是url
	OnLocation(user string, lat, lon, precision float64) ReplyMessage //Latitude地理位置纬度Longitude	地理位置经度Precision	地理位置精度
}

//事件
type event struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      MsgType  `xml:"MsgType"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
	Ticket       string   `xml:"Ticket"`
	Latitude     float64  `xml:"Latitude"`
	Longitude    float64  `xml:"Longitude"`
	Precision    float64  `xml:"Precision"`
}

type MessageReply struct {
	context     *WxPublicApplication
	eventHandle EventHandle
	msgHandle   MessageHandle
}

func (this *MessageReply) SetHandle(event EventHandle, msg MessageHandle) {
	this.eventHandle = event
	this.msgHandle = msg
}

func (this *MessageReply) Validate(msg WxValidateRequest) string {
	theSign := Signature(this.context.Token, msg.Nonce, msg.Timestamp)
	if msg.Signature == theSign {
		return msg.Echostr
	}

	return ""

}

func (this *MessageReply) Dispatch(msg MixMessage) interface{} {
	msgType := msg.MsgType

	var replyMsg *ReplyMessage
	if msgType == "event" {
		replyMsg = this.dispatchEvents(msg)
	} else {
		replyMsg = this.dispatchMessage(msg)
	}

	//转换构造wx消息体对象
	if replyMsg != nil {
		return this.convert(msg, replyMsg)

	}

	return nil
}

func copyFromSource(source MixMessage, target *CommonToken) {
	target.ToUserName = source.FromUserName
	target.FromUserName = source.ToUserName
	target.CreateTime = GetCurrTimeStamps()

}

func (this *MessageReply) convert(source MixMessage, msg *ReplyMessage) interface{} {
	if msg != nil {
		switch msg.Type {
		case MSG_Text:
			targetMsg := Text{}
			targetMsg.MsgType = MsgTypeText
			copyFromSource(source, &targetMsg.CommonToken)
			targetMsg.Content = msg.Content
			return targetMsg
		case MSG_Image:
			targetMsg := Image{}
			targetMsg.MsgType = MsgTypeImage
			copyFromSource(source, &targetMsg.CommonToken)
			targetMsg.Image.MediaID = msg.MediaId
			return targetMsg

		case MSG_Voice:
		case MSG_Video:

		}
	}
	return nil

}

func (this *MessageReply) dispatchMessage(msg MixMessage) *ReplyMessage {
	return nil
}

func (this *MessageReply) dispatchEvents(msg MixMessage) *ReplyMessage {
	if this.eventHandle != nil {

	}

	return nil

}

func (this *MessageReply) DispatchEncrypted(msg WxValidateRequest, xmlmsg EncryptedXMLMsg) interface{} {
	theSign := Signature(this.context.Token, msg.Nonce, msg.Timestamp)
	if msg.Signature == theSign {
		_, rawXMLMsgBytes, err := DecryptMsg(this.context.AppId, xmlmsg.EncryptedMsg, this.context.EncodingAESKey)

		if err == nil {
			theMsg := MixMessage{}
			xml.Unmarshal(rawXMLMsgBytes, &theMsg)
			response := this.Dispatch(theMsg)
			if response != nil {

				//构造加密消息
				responseMsg, _ := xml.Marshal(response)
				var encryptedMsg []byte
				encryptedMsg, err = EncryptMsg([]byte(RandomStr(16)), responseMsg, this.context.AppId, this.context.EncodingAESKey)
				if err != nil {
					return nil
				}

				//生成签名
				timestamp := GetCurrTimeStamps()
				timestampStr := strconv.FormatInt(timestamp, 10)
				theNonce := msg.Nonce //使用微信传递的nonce
				msgSignature := Signature(this.context.Token, timestampStr, theNonce, string(encryptedMsg))

				return &ResponseEncryptedXMLMsg{
					EncryptedMsg: string(encryptedMsg),
					MsgSignature: msgSignature,
					Timestamp:    timestamp,
					Nonce:        theNonce,
				}

			}

		}
	}
	return nil
}
