package utils

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(subject string, body string, to string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("SENDER_MAIL"))

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SENDER_MAIL"), os.Getenv("SENDER_PASS"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
