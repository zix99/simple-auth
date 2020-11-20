package engine

import (
	"net/smtp"
	"strings"
)

type smtpEngine struct {
	host     string
	identity string
	username string
	password string
}

func extractHostname(host string) string {
	idx := strings.Index(host, ":")
	if idx < 0 {
		return host
	}
	return host[:idx]
}

func NewSMTPEngine(host, identity, username, password string) EmailEngine {
	return &smtpEngine{
		host,
		identity,
		username,
		password,
	}
}

func (s *smtpEngine) Send(to, from string, data []byte) error {
	auth := smtp.PlainAuth(s.identity, s.username, s.password, extractHostname(s.host))
	return smtp.SendMail(s.host, auth, from, []string{to}, data)
}
