package mail

import (
	"crypto/tls"
	"fmt"

	smtplib "gopkg.in/mail.v2"
)

type QQmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewQQmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &QQmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *QQmailSender) SendEmail(subject, body string, to, cc, bcc, attachFiles []string) error {
	
	// 发送邮件
	m := smtplib.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress))  
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, f := range attachFiles {
		m.Attach(f)
	}

	d := smtplib.NewDialer("smtp.qq.com", 465, sender.fromEmailAddress, sender.fromEmailPassword)
	
	d.StartTLSPolicy = smtplib.MandatoryStartTLS
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	
	return d.DialAndSend(m)
}


