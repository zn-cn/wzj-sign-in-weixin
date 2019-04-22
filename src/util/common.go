package util

import (
	"config"

	gomail "gopkg.in/gomail.v2"
)

func SendEmail(name, subject, content string, emailTos []string) {
	m := gomail.NewMessage()
	emailInfo := config.Conf.EmailInfo
	m.SetAddressHeader("From", emailInfo.From, name) // 发件人

	// 收件人
	m.SetHeader("To",
		emailTos...,
	)
	m.SetHeader("Subject", subject) // 主题
	m.SetBody("text/html", content) // 正文

	d := gomail.NewPlainDialer(emailInfo.Host, 465, emailInfo.From, emailInfo.AuthCode) // 发送邮件服务器、端口、发件人账号、发件人密码
	d.DialAndSend(m)
}
