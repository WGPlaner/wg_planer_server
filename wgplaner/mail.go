package wgplaner

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
)

var mailLog = logging.MustGetLogger("Mail")

func ValidateMailConfig(config mailConfig) ErrorList {
	err := ErrorList{}

	if !IntInSlice(config.SMTPPort, []int{25, 465, 587}) {
		mailLog.Warning("SMTP Port is not a default port!")
	}

	return err
}

func SendMail(to []string, subject string, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		AppConfig.Mail.SMTPUser,
		AppConfig.Mail.SMTPPassword,
		AppConfig.Mail.SMTPHost,
	)

	mailLog.Debug("Start sending mail")

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
		AppConfig.Mail.SMTPHost+":"+strconv.Itoa(AppConfig.Mail.SMTPPort),
		auth,
		AppConfig.Mail.SMTPIdentity,
		to,
		message,
	)

	if err == nil {
		mailLog.Info("Successfully sent mail!")
	}

	return err
}

// Send a test mail to check if the SMTP connection works!
func SendTestMail() {
	err := SendMail(
		[]string{AppConfig.Mail.SMTPIdentity},
		"WGPlaner Server works",
		"If you get this mail, it means that the server was started successfully!",
	)
	if err != nil {
		mailLog.Fatal("Sending Test Mail failed! Configure the SMTP server! ", err)
	}
}
