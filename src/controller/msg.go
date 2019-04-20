package controller

import (
	"constant"
	"errors"
	"fmt"
	"model"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var event = map[string]model.EventFunc{
	constant.EventTypeSubscribe:   model.EventSubscribe,
	constant.EventTypeUnsubscribe: model.EventUnSubscribe,
	constant.EventTypeLocation:    model.EventLocation,
	constant.EventTypeClick:       model.EventClick,
	constant.EventTypeScan:        model.EventScan,
	constant.EventTypeVIEW:        model.EventView,
}

var stateMachine = map[int]func(string, int) (string, error){
	0: delStateMachine,
	1: delStateMachine,
	2: delStateMachine,
	3: delStateMachine,
	4: func(openid string, status int) (string, error) {
		return getUserCoordinates(openid)
	},
}

var textMachine = map[string]func(string, string) (string, error){
	"帮助":   getHelpText,
	"help": getHelpText,
}

func getHelpText(openid, text string) (string, error) {
	return constant.WechatHelpText, nil
}

func delStateMachine(openid string, status int) (string, error) {
	err := model.SetRedisUserStatus(openid, status)
	if err != nil {
		return "", err
	}
	return constant.WechatTexts[status], nil
}

func getUserCoordinates(openid string) (string, error) {
	coordinates, err := model.GetUserCoordinates(openid)
	if err != nil {
		return "", err
	}

	text := ""
	for tag, coordinate := range coordinates {
		text += fmt.Sprintf(constant.WechatCoordinateFormat+"\n", tag, coordinate.Lon, coordinate.Lat)
	}
	return text, nil
}

var stateDel = map[int]func(string, string) (string, error){
	1: addSignInTask,
	2: addUserCoordinate,
	3: setUserCurCoordinateByTag,
}

// 直接复制链接或者输入openid均可
func addSignInTask(openid, text string) (string, error) {
	u, err1 := url.Parse(text)
	m, err2 := url.ParseQuery(u.RawQuery)
	var err error
	if err1 == nil && err2 == nil {
		if id, ok := m["openid"]; ok {
			err = model.AddSignInTask(openid, id[0])
		} else {
			err = errors.New("empty")
		}
	} else {
		err = errors.New("parse error")
	}

	if err != nil {
		err = model.AddSignInTask(openid, text)
	}
	return constant.WechatAddSignInSuccess, err
}

func addUserCoordinate(openid, text string) (string, error) {
	var tag string
	coordinate := model.Coordinate{}
	strs := strings.Split(text, ":")
	if len(strs) != 2 {
		// 中文
		strs = strings.Split(text, "：")
		if len(strs) != 2 {
			return "", errors.New("格式错误")
		}
	}
	tag = strs[0]

	cs := strings.Split(strs[1], ",")
	if len(cs) != 2 {
		// 中文
		cs = strings.Split(strs[1], "，")
		if len(cs) != 2 {
			return "", errors.New("格式错误")
		}
	}

	var err error
	coordinate.Lon, err = strconv.ParseFloat(cs[0], 64)
	coordinate.Lat, err = strconv.ParseFloat(cs[1], 64)

	err = model.AddUserCoordinate(openid, tag, coordinate)
	return "添加成功", err
}

func setUserCurCoordinateByTag(openid, text string) (string, error) {
	coordinate, err := model.SetUserCurCoordinateByTag(openid, text)
	return fmt.Sprintf("设置成功\n"+constant.WechatCoordinateFormat, text, coordinate.Lon, coordinate.Lat), err
}

// DelEvent 处理微信的消息事件
func DelEvent(c echo.Context) error {
	data := model.EventMsg{}
	err := c.Bind(&data)
	if err != nil {
		return c.XML(http.StatusBadGateway, nil)
	}
	var resData interface{}
	textData := model.WechatTextRes{
		ToUserName:   data.FromUserName,
		FromUserName: data.ToUserName,
		CreateTime:   data.CreateTime,
		MsgType:      "text",
	}
	if data.MsgType == "event" {
		// 处理事件消息
		resData, err = event[data.EventType](data)
	} else {
		// 接收普通消息
		text := data.Content
		status, err := strconv.Atoi(text)
		if f, ok := stateMachine[status]; err == nil && ok {
			// 设置状态
			textData.Content, err = f(data.FromUserName, status)
		} else {
			// 处理状态下的功能
			status, err = model.GetRedisUserStatus(data.FromUserName)
			if delF, ok := stateDel[status]; ok {
				textData.Content, err = delF(data.FromUserName, data.Content)
				// 复原状态
				if err == nil {
					go model.SetRedisUserStatus(data.FromUserName, 0)
				}
			} else {
				if textDelF, ok := textMachine[text]; ok {
					textData.Content, err = textDelF(data.FromUserName, text)
				} else {
					err = errors.New("没有开发此功能")
				}
			}
		}
		resData = textData
	}

	if err != nil {
		resData = model.GetDefaultTextMsg(data)
	}
	return c.XML(http.StatusOK, &resData)
}

type SignatureCheck struct {
	Signature string `json:"signature" query:"signature" form:"signature" validate:"required"`
	Timestamp string `json:"timestamp" query:"timestamp" form:"timestamp" validate:"required"`
	Nonce     string `json:"nonce" query:"nonce" form:"nonce" validate:"required"`
	Echostr   string `json:"echostr" query:"echostr" form:"echostr" validate:"required"`
}

// DelSignature 处理微信服务器的回调URL确认
func DelSignature(c echo.Context) error {
	data := SignatureCheck{}
	err := c.Bind(&data)
	echostr := ""
	if err != nil {
		return c.String(http.StatusOK, echostr)
	}
	signature := model.GetSignature(data.Timestamp, data.Nonce)
	if signature == data.Signature {
		echostr = data.Echostr
	}
	// oh, fork, 不能返回json
	return c.String(http.StatusOK, echostr)
}
