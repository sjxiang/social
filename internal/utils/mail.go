package utils

import (
	"bytes"
	"text/template"
)

// 邮件模板参数
type MailTemplateParams struct {
	Username      string
	ActivationURL string
}

func BuildPlainTextMessage(params MailTemplateParams) (string, error) {
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

	// 定义模板, 解析模板
	tpl, err := template.New("user_invitation").Parse(example)
	if  err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}

	// 渲染模板
	if err := tpl.Execute(buffer, params); err != nil {
		return "", err
	}
	
	formattedMessage := buffer.String()

	return formattedMessage, nil 
}


