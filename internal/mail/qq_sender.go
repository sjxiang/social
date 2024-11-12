package mail

import (
	"bytes"
	"text/template"
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

func (sender *QQmailSender) SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error {
	
	// 发送邮件
	m := smtplib.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress))  
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

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


type Params struct {
	Username      string
	ActivationURL string
}

func FormattedContent(arg Params) (string, error) {
	example := `
您好，{{ .Username }}<br>
<br>
请点击下面的URL，以完成「xxx 开发者社区」的账户认证。<br>
<br>
&nbsp;∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴<br>
<br>
<a href="{{ .ActivationURL }}">https://clwy.cn/users/confirmation?confirmation_token=9D2K1sGA_x2dmVj9npF4</a>
<br>
<br>
&nbsp;∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴ ∴<br>
<br>
※<wbr>如果您无法点击以上URL<wbr><br>
麻烦您复制上面的URL，<wbr>粘贴到浏览器的地址栏中访问。<br>
<br>
※如果您不知道为什么收到此邮件。<wbr>非常抱歉，请忽略这封邮件。<br>
※这封邮件是认证账户专用的。<wbr>请不要回复此邮件，请谅解。<br>
<br>
━━━━━━━━━━━━━━━━<br>
社区小管家<br>`

	// 定义模板
	tpl := template.New("user_invitation")

	// 解析模板
	tpl, err := tpl.Parse(example)
	if  err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}

	// 渲染模板
	if err := tpl.Execute(buffer, arg); err != nil {
		return "", err
	}
	
	return buffer.String(), nil 
}


