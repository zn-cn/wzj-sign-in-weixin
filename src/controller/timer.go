package controller

import "model"

func StartHourTimer() {
	// 先更新 access_token 后更新ticket
	model.UpdateAccessToken()
}

func StartTaskTimer() {
	model.StartStuSignTask()
}
