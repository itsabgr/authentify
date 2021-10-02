package authentify

import (
	"errors"
	"fmt"
	"net/smtp"
	"time"
)

type emailSender struct {
	from, server string
	auth         smtp.Auth
}

func (s *emailSender) Send(to string, code Code, exp time.Time) error {
	if exp.Before(time.Now()) {
		return errors.New("expired")
	}
	msg := []byte(fmt.Sprintf("Code: %s\nExpire At: %s\n", code.String(), exp))
	return smtp.SendMail(s.server, s.auth, s.from, []string{to}, msg)
}

func (s *emailSender) Name() string {
	return "email"
}

func NewEmailSender(server, from string, auth smtp.Auth) (Sender, error) {
	sender := &emailSender{}
	sender.auth = auth
	sender.server = server
	sender.from = from
	return sender, nil
}
