package main


// "Invalid credentials."


type Message struct {
	To       string
	Subject  string
	Data     []string
	Template string
}
// 10:25
// 11:14 19天，3个男生

// 14
// sendEmail()

type Text struct {
	Body string

}


func SendMail() {
	// 1. 从redis中获取邮件信息
	// 2. 发送邮件
	// 3. 删除邮件信息
}
// msg := Message{
// 	To:       user.Email,
// 	Subject:  "验证电子邮箱",
// 	Data:     invoice,
// 	Template: "invoice",
// }
