package controller

import "github.com/yun-mu/wzj-sign-in-weixin/model"

func StartHourTimer() {
	// 先更新 access_token 后更新ticket
	model.UpdateAccessToken()
}

func StartTaskTimer() {
	model.StartStuSignTask()
}
