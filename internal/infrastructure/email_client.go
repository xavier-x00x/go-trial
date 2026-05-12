package infrastructure

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailClient(host string, port int, username, password, fromEmail, fromName string) *EmailClient {
	dialer := &gomail.Dialer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	return &EmailClient{
		dialer: dialer,
		from:   fmt.Sprintf("%s <%s>", fromName, fromEmail),
	}
}

func (c *EmailClient) SendEmail(to, subject, body string) error {
	return c.SendEmailWithAttachment(to, subject, body, nil, "")
}

func (c *EmailClient) SendEmailWithAttachment(
	to, subject, body string,
	attachment []byte, filename string,
) error {
	if c.dialer.Username == "" || c.from == "" {
		return fmt.Errorf("SMTP not configured")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", c.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	if len(attachment) > 0 && filename != "" {
		tmpFile, err := os.CreateTemp("", filename)
		if err == nil {
			defer os.Remove(tmpFile.Name())
			tmpFile.Write(attachment)
			tmpFile.Close()
			msg.Attach(tmpFile.Name())
		}
	}

	if err := c.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}