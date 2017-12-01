package wx

import (
	"fmt"
	"time"
)

const (
	API_GETWXIPS_URL = "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s" //获取微信服务器ip列表
	//获取accesstoken
	API_GETTOKEN_URL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

//微信公众号应用
type WxPublicApplication struct {
	AppId          string    //应用id
	AppSecret      string    //应用秘钥
	Token          string    //访问设定的token
	EncodingAESKey string    //AES秘钥
	UpdateTime     time.Time //token的更新时间
	NextUpdateTime time.Time //token下次更新时间
	token          *WxAccessToken
}

type WxAccessToken struct {
	CommonError
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires_in"`
}

//GetAccessToken 获取access_token，如果accesstoken不存在或已经失效则会重新获取
func (this *WxPublicApplication) GetAccessToken() (accessToken string, err error) {
	if this.token == nil {
		err = this.getAccessTokenFromWx()

	} else {
		theNow := time.Now()
		if theNow.After(this.NextUpdateTime) {
			err = this.getAccessTokenFromWx()
		}
	}

	if this.token != nil {
		accessToken = this.token.AccessToken
	}

	return
}

func (this *WxPublicApplication) getAccessTokenFromWx() (err error) {
	result := new(WxAccessToken)
	err = WxGet(API_GETTOKEN_URL, this.AppId, this.AppSecret)
	if err == nil {
		if result.ErrCode == WX_SUCCESS {
			this.UpdateTime = time.Now()
			this.NextUpdateTime = this.UpdateTime.Add(time.Duration(result.Expires) * time.Second)
			this.token = result
			return nil
		}
		return buildErrorByCode(result.ErrCode)

	}
	return err
}

//微信服务器列表
type WxServiceIps struct {
	List []string `json:"ip_list"`
}

//获取微信服务器列表
func (this *WxPublicApplication) GetWxServiceIps() (list WxServiceIps, err error) {
	result := new(WxServiceIps)
	token, err := this.GetAccessToken()
	if err == nil {
		err = WxGet(API_GETWXIPS_URL, result, token)
		return *result, err
	}

	return *result, nil
}

func buildErrorByCode(code int64) (err error) {
	switch code {

	case -1:
		err = fmt.Errorf(ERROR_TEMPLATE, code, "系统繁忙，此时请开发者稍候再试")
	case 40001:
		err = fmt.Errorf(ERROR_TEMPLATE, code, "AppSecret错误或者AppSecret不属于这个公众号，请开发者确认AppSecret的正确性")
	case 40002:
		err = fmt.Errorf(ERROR_TEMPLATE, code, "请确保grant_type字段值为client_credential")
	case 40164:
		err = fmt.Errorf(ERROR_TEMPLATE, code, "调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置")

	}
	return

}
