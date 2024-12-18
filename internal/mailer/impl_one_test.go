package mailer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjxiang/social/internal/utils"
)

var (
	// 发件人邮箱账号
	sender = "gua@vip.cn"
	// 发件人邮箱密码，这里需要注意，QQ邮箱如果开启了SMTP服务，可能需要使用授权码而不是邮箱登录密码
	password = "xxx"
)

func TestSendMail(t *testing.T) {

	// '-short' tag, 会跳过该项测试 
	if testing.Short() {
		t.Skip()
	}

	// 收件人邮箱账号
	receiver := "sjxiang2020@outlook.com"

	// 测试邮件发送
	sender := NewQQmailSender("no-reply", sender, password)
	
	params := utils.MailTemplateParams{
		Username:      "xxx",
		ActivationURL: "www.baidu.com",
	}
	content, err := utils.BuildPlainTextMessage(params)
	require.NoError(t, err)

	if err := sender.SendEmail("这是一封测试邮件", content, []string{receiver}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}

	t.Log("邮件发送成功")
}

