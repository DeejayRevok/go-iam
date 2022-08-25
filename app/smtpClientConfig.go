package app

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func BuildSMTPClient() *mail.SMTPClient {
	host := os.Getenv("SMTP_SERVER")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(fmt.Sprintf("Malformed SMTP port: %s", err.Error()))
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	smtpServer := mail.NewSMTPClient()
	smtpServer.Host = host
	smtpServer.Port = port
	smtpServer.Username = user
	smtpServer.Password = password
	smtpServer.Encryption = mail.EncryptionTLS

	smtpServer.KeepAlive = true
	smtpServer.ConnectTimeout = 10 * time.Second
	smtpServer.SendTimeout = 10 * time.Second
	smtpServer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := smtpServer.Connect()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the SMTP server: %s", err.Error()))
	}
	return smtpClient
}
