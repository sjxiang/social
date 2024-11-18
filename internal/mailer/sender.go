package mailer


type EmailSender interface {
	SendEmail(subject, body string, to, cc, bcc, attachFiles []string) error
}


/*

参数
	subject 主题
	body 正文
	to 收件人, 群发
	cc 抄送, 可见(套瓷, 尴尬)
	bcc 抄送, 不可见
	attachFiles 附件

	
 */