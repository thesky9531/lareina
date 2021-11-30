package email

import (
	"gopkg.in/gomail.v2"
)

type GoMailConfig struct {
	Title         string   // 邮件标题
	Body          string   // 邮件内容
	RecipientList []string // 收件人列表
	Sender        string   // 发件人账号
	SPassword     string   // 发件人密码，QQ邮箱这里配置授权码
	SMTPAddr      string   // SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPPort      int      // SMTP端口 QQ邮箱是25
}

func NewGoMailEntity() (m GoMailConfig) {
	return
}

func (m *GoMailConfig) SetGoMailEntity(username, password, host, title, body string, port int, to []string) {
	m.Title = title
	m.Body = body
	m.RecipientList = to
	m.Sender = username
	m.SPassword = password
	m.SMTPAddr = host
	m.SMTPPort = port
}
func (m *GoMailConfig) SetGoMailBody(body string) {
	m.Body = body
}

func (m GoMailConfig) SendMail() error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.Sender)
	mail.SetHeader("To", m.RecipientList...)
	mail.SetHeader("Subject", m.Title)
	mail.SetBody(`text/html`, m.Body)
	err := gomail.NewDialer(m.SMTPAddr, m.SMTPPort, m.Sender, m.SPassword).DialAndSend(mail)
	return err
}
