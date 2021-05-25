package email

import (
	"net/smtp"
	"strings"
)

type Mail interface {
}

type mailEntity struct {
	user     string
	password string
	host     string
	to       string
	subject  string
	body     string
	mailtype string
}

func NewMailEntity() (m mailEntity) {
	return m
}
func (m mailEntity) SetMailEntity(user, password, host, to, subject, body, mailtype string) {
	m.user = user
	m.password = password
	m.host = host
	m.to = to
	m.subject = subject
	m.body = body
	m.mailtype = mailtype
}

func (m mailEntity) SendToMail() error {
	hp := strings.Split(m.host, ":")
	auth := smtp.PlainAuth("", m.user, m.password, hp[0])
	content_type := "Content-Type: text/"
	if m.mailtype == "html" {
		content_type += m.mailtype
		content_type += "; charset=UTF-8"
	} else {
		content_type += "plain; charset=UTF-8"
	}
	msg := []byte("To: " + m.to + "\r\nFrom: " + m.user + "\r\nSubject: " + m.subject + "\r\ncontent_type: " + content_type + "\r\n\r\n" + m.body)
	sendTo := strings.Split(m.to, ";")
	err := smtp.SendMail(m.host, auth, m.user, sendTo, msg)
	return err
}
