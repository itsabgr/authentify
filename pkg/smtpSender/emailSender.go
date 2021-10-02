package smtpSender

import (
	"context"
	"errors"
	"github.com/itsabgr/authentify"
	"net/smtp"
	"time"
)

type emailSender struct {
	from, server string
	auth         smtp.Auth
	msgRenderer  MsgRenderer
}

func (s *emailSender) Close() error {
	s.server = ""
	s.from = ""
	s.auth = nil
	s.msgRenderer = nil
	return nil
}

func (s *emailSender) Send(ctx context.Context, to string, code authentify.Code, exp time.Time) error {
	if exp.Before(time.Now()) {
		return errors.New("expired")
	}
	msg, err := s.msgRenderer.RenderMsg(to, code, exp)
	if err != nil {
		return err
	}
	return smtp.SendMail(s.server, s.auth, s.from, []string{to}, msg)
}

func (s *emailSender) Name() string {
	return "email"
}

type MsgRenderer interface {
	RenderMsg(to string, code authentify.Code, exp time.Time) ([]byte, error)
}

func NewEmailSender(server, from string, mr MsgRenderer, auth smtp.Auth) (authentify.Sender, error) {
	sender := &emailSender{}
	sender.auth = auth
	sender.server = server
	sender.from = from
	sender.msgRenderer = mr
	return sender, nil
}
