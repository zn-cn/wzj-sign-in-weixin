/*
Package main package is the entry file
*/
package main

import (
	"github.com/yun-mu/wzj-sign-in-weixin/controller"

	"github.com/yun-mu/wzj-sign-in-weixin/config"
	"github.com/yun-mu/wzj-sign-in-weixin/constant"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/robfig/cron"
)

func main() {
	go startTimer()
	startWeb()
}

func startWeb() {
	e := echo.New()

	e.Use(middleware.Recover())

	v1 := e.Group(constant.APIPrefix)
	v1.GET("/event", controller.DelSignature)
	v1.POST("/event", controller.DelEvent)

	e.Logger.Fatal(e.Start(config.Conf.AppInfo.Addr))
}

func startTimer() {
	c := cron.New()
	controller.StartHourTimer()
	c.AddFunc(constant.TimerEveryHour, controller.StartHourTimer)
	c.AddFunc(constant.TimerEveryTwentySecond, controller.StartTaskTimer)
	c.Start()
}
