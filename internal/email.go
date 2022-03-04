package internal

import (
	"net/smtp"
)

func SendConfirmationEmail(to []string, id string) bool {
	from := "lipearthur81@gmail.com"
	pmail := "egcdbdpzhljhgnkh"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("Subject: E-mail de confirmação de conta\r\nContent-Type: text/plain; charset=\"utf-8\"\r\nConfirmation e-mail\nhttp://localhost:1337/user/confirm/" + id)

	auth := smtp.PlainAuth("", from, pmail, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg); err != nil {
		return false
	}

	return true
}
