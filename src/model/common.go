package model

import (
	"config"
	"constant"
	"fmt"
	"sync"

	"github.com/imroc/req"
)

type weixinAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

var (
	appid     string
	appSecret string

	accessToken string
	tokenMutex  sync.RWMutex
)

func init() {
	appid = config.Conf.Wechat.AppID
	appSecret = config.Conf.Wechat.AppSecret
}

func GetDefaultTextMsg(msg EventMsg) WechatTextRes {
	resTextMsg := NewResTextMsg(msg)
	resTextMsg.Content = constant.WechatDefaultText
	return resTextMsg
}

func GetWelcomeTextMsg(msg EventMsg) WechatTextRes {
	resTextMsg := NewResTextMsg(msg)
	resTextMsg.Content = constant.WechatWelcomeText
	return resTextMsg
}

func NewResTextMsg(msg EventMsg) WechatTextRes {
	return WechatTextRes{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   msg.CreateTime,
		MsgType:      "text",
	}
}

func requestAccessToken() (weixinAccessTokenResponse, error) {
	url := fmt.Sprintf(constant.URLWechatAccessToken, appid, appSecret)
	r, err := req.Get(url)
	var respBody weixinAccessTokenResponse
	if err != nil {
		return respBody, err
	}
	err = r.ToJSON(&respBody)
	return respBody, err
}

func UpdateAccessToken() error {
	resp, err := requestAccessToken()
	if err != nil {
		return err
	}

	if resp.AccessToken == "" {
		return fmt.Errorf("errcode: %d, errmsg: %s", resp.Errcode, resp.Errmsg)
	}
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	accessToken = resp.AccessToken
	return nil
}
