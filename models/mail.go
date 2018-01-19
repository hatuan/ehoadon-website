package models

import (
	"bytes"
	"html/template"
	"net/smtp"
)

var auth smtp.Auth

type Mail struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewMail(to []string, subject, body string) *Mail {
	auth = smtp.PlainAuth("", "erpvietnam.dyndns@gmail.com", "pass@wOrd1", "smtp.gmail.com")

	return &Mail{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Mail) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "erpvietnam.dyndns@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Mail) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
