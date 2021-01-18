package smtp

import (
	"errors"
	"fmt"
	"net/smtp"
	"strconv"
)

type SMTPEmailConfig struct {
	Email    string
	Password string
	Hostname string
	Port     uint
}

type SMTPEmail struct {
	config SMTPEmailConfig
	auth   smtp.Auth
}

func (s SMTPEmail) Send(to, title, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Mime-Version: 1.0;\r\n"+
		"Content-Type: text/html;\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", title, to, body))
	return smtp.SendMail(s.config.Hostname+":"+strconv.Itoa(int(s.config.Port)), s.auth, s.config.Email, []string{to}, msg)
}

func NewEmail(config SMTPEmailConfig) (SMTPEmail, error) {
	if config.Email == "" || config.Password == "" || config.Hostname == "" || config.Port == 0 {
		return SMTPEmail{}, errors.New("smtp: config fields can't be empty")
	}

	return SMTPEmail{
		config: config,
		auth:   smtp.PlainAuth("", config.Email, config.Password, config.Hostname),
	}, nil
}
