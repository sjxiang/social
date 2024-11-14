package mail


type EmailSender interface {
	SendEmail(
		subject string,        // 邮件主题
		body string,           // 邮件正文
		to []string,           // 收件人
		cc []string,           // 抄送, 可见(套瓷, 尴尬)
		bcc []string,          // 抄送, 不可见 
		attachFiles []string,  // 附件
	) error
}

