package wx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	redirectOauthURL      = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	accessTokenURL        = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	the_userInfo_URL      = "https://api.weixin.qq.com/cgi-bin/user/info"
	checkAccessTokenURL   = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

//Oauth 保存用户授权信息
type Oauth struct {
	context *WxPublicApplication
}

//GetRedirectURL 获取跳转的url地址
func (oauth *Oauth) GetRedirectURL(redirectURI, scope, state string) (string, error) {
	//url encode
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(redirectOauthURL, oauth.context.AppId, urlStr, scope, state), nil
}

//Redirect 跳转到网页授权
func (oauth *Oauth) Redirect(writer http.ResponseWriter, redirectURI, scope, state string) error {
	location, err := oauth.GetRedirectURL(redirectURI, scope, state)
	if err != nil {
		return err
	}
	//location 为完整地址，所以不需要request
	http.Redirect(writer, nil, location, 302)
	return nil
}

// ResAccessToken 获取用户授权access_token的返回结果
type ResAccessToken struct {
	CommonError

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (oauth *Oauth) GetUserAccessToken(code string) (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, oauth.context.AppId, oauth.context.AppSecret, code)
	var response []byte
	response, err = HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

//RefreshAccessToken 刷新access_token
func (oauth *Oauth) RefreshAccessToken(refreshToken string) (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(refreshAccessTokenURL, oauth.context.AppId, refreshToken)
	var response []byte
	response, err = HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

//CheckAccessToken 检验access_token是否有效
func (oauth *Oauth) CheckAccessToken(accessToken, openID string) (b bool, err error) {
	urlStr := fmt.Sprintf(checkAccessTokenURL, accessToken, openID)
	var response []byte
	response, err = HTTPGet(urlStr)
	if err != nil {
		return
	}
	var result CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		b = false
		return
	}
	b = true
	return
}

//UserInfo 用户授权获取到用户信息
type UserInfo struct {
	CommonError

	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

//Info 用户基本信息
type User_Info struct {
	UserInfo
	Subscribe     int32    `json:"subscribe"`
	Language      string   `json:"language"`
	SubscribeTime int32    `json:"subscribe_time"`
	Remark        string   `json:"remark"`
	GroupID       int32    `json:"groupid"`
	TagidList     []string `json:"tagid_list"`
}

//GetUserInfo 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (oauth *Oauth) GetUserInfo(accessToken, openID string) (result UserInfo, err error) {
	urlStr := fmt.Sprintf(userInfoURL, accessToken, openID)
	var response []byte
	response, err = HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

//GetUserInfo 获取用户基本信息
func (this *Oauth) GetUser_Info(openID string) (userInfo *User_Info, err error) {
	var accessToken string
	accessToken, err = this.context.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&openid=%s&lang=zh_CN", the_userInfo_URL, accessToken, openID)
	var response []byte
	response, err = HTTPGet(uri)
	if err != nil {
		return
	}
	userInfo = new(User_Info)
	err = json.Unmarshal(response, userInfo)
	if err != nil {
		return
	}
	if userInfo.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo Error , errcode=%d , errmsg=%s", userInfo.ErrCode, userInfo.ErrMsg)
		return
	}
	return
}
