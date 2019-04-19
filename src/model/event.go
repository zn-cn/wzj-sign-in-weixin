package model

import (
	"config"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"sort"
)

// 微信服务器推送过来的消息(事件)的通用消息头.
type MsgHeader struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string `xml:"MsgType"      json:"MsgType"`
}

// 微信服务器推送过来的消息(事件)的合集.
type EventMsg struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MsgHeader
	EventType string `xml:"Event" json:"Event"`

	MsgId        int64   `xml:"MsgId"        json:"MsgId"`
	Content      string  `xml:"Content"      json:"Content"`
	MediaId      string  `xml:"MediaId"      json:"MediaId"`
	PicURL       string  `xml:"PicUrl"       json:"PicUrl"`
	Format       string  `xml:"Format"       json:"Format"`
	Recognition  string  `xml:"Recognition"  json:"Recognition"`
	ThumbMediaId string  `xml:"ThumbMediaId" json:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"   json:"Location_X"`
	LocationY    float64 `xml:"Location_Y"   json:"Location_Y"`
	Scale        int     `xml:"Scale"        json:"Scale"`
	Label        string  `xml:"Label"        json:"Label"`
	Title        string  `xml:"Title"        json:"Title"`
	Description  string  `xml:"Description"  json:"Description"`
	URL          string  `xml:"Url"          json:"Url"`
	EventKey     string  `xml:"EventKey"     json:"EventKey"`
	Ticket       string  `xml:"Ticket"       json:"Ticket"`
	Latitude     float64 `xml:"Latitude"     json:"Latitude"`
	Longitude    float64 `xml:"Longitude"    json:"Longitude"`
	Precision    float64 `xml:"Precision"    json:"Precision"`

	// menu
	MenuId       int64 `xml:"MenuId" json:"MenuId"`
	ScanCodeInfo *struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`
		ScanResult string `xml:"ScanResult" json:"ScanResult"`
	} `xml:"ScanCodeInfo,omitempty" json:"ScanCodeInfo,omitempty"`
	SendPicsInfo *struct {
		Count   int `xml:"Count" json:"Count"`
		PicList []struct {
			PicMd5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"`
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"`
	} `xml:"SendPicsInfo,omitempty" json:"SendPicsInfo,omitempty"`
	SendLocationInfo *struct {
		LocationX float64 `xml:"Location_X" json:"Location_X"`
		LocationY float64 `xml:"Location_Y" json:"Location_Y"`
		Scale     int     `xml:"Scale"      json:"Scale"`
		Label     string  `xml:"Label"      json:"Label"`
		PoiName   string  `xml:"Poiname"    json:"Poiname"`
	} `xml:"SendLocationInfo,omitempty" json:"SendLocationInfo,omitempty"`

	MsgID  int64  `xml:"MsgID"  json:"MsgID"`  // template, mass
	Status string `xml:"Status" json:"Status"` // template, mass
}

type WechatTextRes struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string   `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string   `xml:"MsgType"      json:"MsgType"`
	Content      string   `xml:"Content" json:"Content"`
}

type EventFunc func(EventMsg) (interface{}, error)

func EventSubscribe(msg EventMsg) (interface{}, error) {
	err := createUser(msg.FromUserName)
	return GetWelcomeTextMsg, err
}

func EventUnSubscribe(msg EventMsg) (interface{}, error) {
	return GetDefaultTextMsg, nil
}

func EventScan(msg EventMsg) (interface{}, error) {
	return GetDefaultTextMsg, nil
}

func EventLocation(msg EventMsg) (interface{}, error) {
	return GetDefaultTextMsg, nil
}

func EventClick(msg EventMsg) (interface{}, error) {
	return GetDefaultTextMsg, nil
}

func EventView(msg EventMsg) (interface{}, error) {
	return GetDefaultTextMsg, nil
}

func MsgEvent(msg EventMsg) (interface{}, error) {
	return GetDefaultTextMsg, nil
}

func GetSignature(timestamp, nonce string) string {
	return GetSign(config.Conf.Wechat.Token, timestamp, nonce)
}

// Sign 微信公众号 url 签名.
func GetSign(token, timestamp, nonce string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce))
	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}
