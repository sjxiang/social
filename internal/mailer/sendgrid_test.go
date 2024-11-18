package mailer
import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendGridSendMail(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	// 测试邮件发送
	sender := NewSendgrid("xxx", "xxx", true)

	statusCode, err := sender.SendEmail("这是一封测试邮件", "gua", "sjok.com", "艹你, 贺惠")
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)
	t.Log("邮件发送成功")
}

