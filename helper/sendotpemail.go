package helper

import (
	"fmt"
	"net/smtp"
	"shoego/config"
)

func SendOTPEmail(email string, otp string) error {

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	from := cfg.EMAIL
	password := cfg.EMAIL_PASSWORD
	smtpHost := cfg.SMTP_HOST
	smtpPort := cfg.SMTP_PORT

	message := []byte("Subject: Your OTP Verification\n\nYour OTP is: " + otp)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort,auth,from,[]string{email},message)

	if err != nil {
		return err
	}

	fmt.Println("OTP email sent successfully")
	return nil
}