package events

import (
	"log"
	"net/smtp"
)

func SendEmailSMTP(to, subject, body string) {
	from := "arailym.mukazhanova2904@gmail.com"
	password := "irciybhianolysqq"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("From: Arailym <" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("[SMTP] ❌ Error: %v", err)
	} else {
		log.Printf("[SMTP] ✅ Sended to: %s", to)
	}
}
