package constant

const (
	APIPrefix = "/api/v1"

	TimerEveryHour       = "@hourly"   // 每小时触发
	TimerEveryFiveSecond = "@every 5s" // 每五秒触发

	/****************************************** wechat ****************************************/
	WechatWelcomeText              = "终于等到你！"
	WechatAddSignInSuccess         = "任务添加成功"
	WechatHelpText                 = "输入指南：\n0 -> 默认状态，输入0也会将状态重置\n\n1 -> 设置为输入链接状态，输入1之后直接输入微助教任意学生栏目下页面复制后的链接（带openid）或者直接输入openid即可添加入签到任务中\n\n注：此openid仅有两小时有效时间，因此每次上课前复制一次在公众号输入即可，每次进入微助教网页openid均会更新，如果重新进入了微助教页面，请更新openid\n\n2 -> 设置为输入坐标状态\n输入格式如下：东十二:114.440465,30.517877\n左侧为标签，右侧为经纬度\n\n3 -> 设置当前坐标标签状态\n输入格式如下：东十二，设置之后之后所有的签到将使用此标签代表的坐标进行签到\n\n4 -> 获取设置的所有标签和坐标\n\n5 -> 获取所有讨论课程\n\n6 -> 设置讨论课程,格式如下:id,课程名\n\n7 -> 设置通知邮件\n\nhelp/帮助 -> 获取帮助指南\n\n坐标拾取：https://lbs.amap.com/console/show/picker"
	WechatDefaultText              = "Sorry，此功能暂未开发。"
	WechatCoordinateFormat         = "标签：%s, 经纬度：%f,%f"
	WechatCourseWithTopicFormat    = "课程名: %s, 课程ID: %d, 话题: %s"
	WechatCourseWithOutTopicFormat = "课程名: %s, 课程ID: %d"

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
	URLWZHDisCourseList  = "https://v18.teachermate.cn/wechat-api/v1/students/courses?type=discussions"
	URLWZJCourseSelect   = "https://v18.teachermate.cn/wechat-api/v1/discussions/select"
)

var (
	WechatTexts = map[int]string{
		0: "重置成功",
		1: "请输入微助教任意页面复制后的链接开启自动签到，注意：输入之后请勿再进入微助教任意页面，否则链接将失效，并需要你重新输入，每个链接有效期为两小时。",
		2: "请输入标签：经度，纬度, eg:东十二:114.440465,30.517877",
		3: "请输入你要设置的当前标签，用于GPS签到，eg:东十二",
	}
)
