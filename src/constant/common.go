package constant

const (
	APIPrefix = "/api/v1"

	TimerEveryHour       = "@hourly"   // 每小时触发
	TimerEveryFiveSecond = "@every 5s" // 每五秒触发

	/****************************************** wechat ****************************************/
	WechatWelcomeText      = "终于等到你！"
	WechatAddSignInSuccess = "任务添加成功"
	WechatHelpText         = "0->重置状态\n1->设置为输入链接状态\n2->设置为输入坐标状态\n3->设置当前坐标标签状态\n4->获取设置的坐标"
	WechatDefaultText      = "Sorry，此功能暂未开发。"
	WechatCoordinateFormat = "标签：%s, 经纬度：%f,%f"

	/****************************************** other ****************************************/

	EventTypeSubscribe   = "subscribe"   // 关注事件, 包括点击关注和扫描二维码(公众号二维码和公众号带参数二维码)关注
	EventTypeUnsubscribe = "unsubscribe" // 取消关注事件
	EventTypeScan        = "SCAN"        // 已经关注的用户扫描带参数二维码事件
	EventTypeLocation    = "LOCATION"    // 上报地理位置事件
	EventTypeClick       = "CLICK"
	EventTypeVIEW        = "VIEW"

	URLWechatAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	URLWZJSignIn         = "https://v18.teachermate.cn/wechat/wechat/guide/signin?openid=%s"
	URLWZJStuSignIn      = "https://v18.teachermate.cn/wechat-api/v1/class-attendance/student-sign-in"
)

var (
	WechatTexts = map[int]string{
		0: "重置成功",
		1: "请输入微助教任意页面复制后的链接开启自动签到，注意：输入之后请勿再进入微助教任意页面，否则链接将失效，并需要你重新输入，每个链接有效期为两小时。",
		2: "请输入标签：经度，纬度, eg:东十二:114.440465,30.517877",
		3: "请输入你要设置的当前标签，用于GPS签到，eg:东十二",
	}
)
