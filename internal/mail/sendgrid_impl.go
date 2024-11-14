package mail

import (
	"fmt"
	"time"

	"github.com/sendgrid/sendgrid-go"
	smtp "github.com/sendgrid/sendgrid-go/helpers/mail"
)


type SendGridMailer struct {
	fromEmail string
	apiKey    string  // 授权码
	isSandbox bool    // 是否开启沙箱模式
	client    *sendgrid.Client  
}

func NewSendgrid(apiKey, fromEmail string, isSandbox bool) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		isSandbox: isSandbox,
		client:    client,
	}
}

func (s *SendGridMailer) SendEmail(subject, username, toEmail, body string) (int, error) {
	// 发件人
	from := smtp.NewEmail("no-reply", s.fromEmail)
	// 收件人
	to := smtp.NewEmail(username, toEmail)
	// 正文
	message := smtp.NewSingleEmail(from, subject, to, "", body)
	
	message.SetMailSettings(&smtp.MailSettings{
		SandboxMode: &smtp.Setting{
			Enable: &s.isSandbox,
		},
	})


	var err error

	for i := 0; i < 3; i++ {

		// 发送邮件
		resp, err := s.client.Send(message)
		
		// 情况一, 多次重试, 寄
		if err != nil {
			// 指数补偿
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		fmt.Println(resp.StatusCode, resp.Body, resp.Headers)

		// 情况二, 成功
		return resp.StatusCode, nil
	}

	return -1, fmt.Errorf("尝试了 3 次, 还是发送失败, %+v", err)
}
