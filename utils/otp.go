package utils

import (
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)



func GenerateOTP (length int)string{
	rand.Seed(time.Now().UnixNano())
	otp := ""
	for i := 0; i < length; i++ {
		otp += strconv.Itoa(rand.Intn(10)) // 0-9
	}
	return otp
}




func SendEmail(to, subject, body string) error {
	from := "aadilnanoth@gmail.com"
	password := "pkfhpnnugjbohane"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
