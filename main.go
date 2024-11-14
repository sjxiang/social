package main

import (
	"log"

	"github.com/sjxiang/social/internal/mail"
)

func main() {

	// 测试邮件发送
	sender := mail.NewSendgrid("xxx", "xxx", true)
	statusCode, err := sender.SendEmail("这是一封测试邮件", "gua", "xxx", "艹傻逼")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(statusCode)

}
