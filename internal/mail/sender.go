package mail


type EmailSender interface {
	SendEmail(
		subject string,        // 邮件主题
		content string,        // 邮件内容
		to []string,           // 收件人邮箱地址
		cc []string,           // 抄送, 可见；套瓷多尴尬
		bcc []string,          // 抄送, 不可见 
		attachFiles []string,  // 附件
	) error
}

