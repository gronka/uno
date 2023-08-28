package uf

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/rs/zerolog/log"
)

type Mailer struct {
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
}

// logMailer is local so only the logger accesses it, but we init the values
// from conf in pile.InitFields
var logMailer Mailer

func (mailer *Mailer) sendLog(subject, body string) {
	headers := make(map[string]string)
	headers["From"] = "taylor@gronka.us"
	headers["Subject"] = "FriLog: " + subject

	switch subject {
	case "PANIC":
		fallthrough
	case "ERROR":
		// for production
		//headers["To"] = "mr.gronka@gmail.com,twmosher@gmail.com"
		headers["To"] = "mr.gronka@gmail.com"

	case "WARN":
		fallthrough
	default:
		headers["To"] = "mr.gronka@gmail.com"
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := mailer.smtpHost + ":" + mailer.smtpPort

	auth := smtp.PlainAuth("", mailer.smtpUser, mailer.smtpPass, mailer.smtpHost)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         servername,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic()
	}

	client, err := smtp.NewClient(conn, mailer.smtpHost)
	if err != nil {
		LoggingError("Mailer: failed to connect")
	}

	if err = client.Auth(auth); err != nil {
		LoggingError("Mailer: failed to auth")
	}

	if err = client.Mail(headers["From"]); err != nil {
		LoggingError("Mailer: mailed to set From")
	}

	if err = client.Rcpt(headers["To"]); err != nil {
		LoggingError("Mailer: mailed to set To")
	}

	writer, err := client.Data()
	if err != nil {
		LoggingError("Mailer: failed to create writer")
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		LoggingError("Mailer: failed to write data")
	}

	err = writer.Close()
	if err != nil {
		LoggingError("Mailer: failed to close writer")
	}

	client.Quit()
}
