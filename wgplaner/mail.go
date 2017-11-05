package wgplaner

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func SendMail(to []string, subject string, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		AppConfig.Mail.SmtpUser,
		AppConfig.Mail.SmtpPassword,
		AppConfig.Mail.SmtpHost,
	)

	log.Println("[Mail] Sending mail")

	message := []byte(
		fmt.Sprintf("Date: %s (UTC)\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			time.Now().UTC().Format(time.RFC1123Z),
			strings.Join(to, ", "),
			subject,
			body,
		),
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		AppConfig.Mail.SmtpHost+":"+strconv.Itoa(AppConfig.Mail.SmtpPort),
		auth,
		AppConfig.Mail.SmtpIdentity,
		to,
		message,
	)

	if err == nil {
		log.Println("[Mail] Successfully sent mail!")
	}

	return err
}

// Send a test mail to check if the SMTP connection works!
func SendTestMail() {
	err := SendMail(
		[]string{AppConfig.Mail.SmtpIdentity},
		"WGPlaner Server works",
		"If you get this mail, it means that the server was started successfully!",
	)
	if err != nil {
		log.Fatalln("[Mail] Sending Test Mail failed! Configure the SMTP server!")
	}
}
