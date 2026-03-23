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

	message := []byte("Subject: Shoego Account Verification OTP\n\n" +
		"Hello,\n\n" +
		"Thank you for using Shoego.\n\n" +
		"Your OTP for verification is: " + otp + "\n\n" +
		"This OTP is valid for 2 minutes.\n\n" +
		"Please do not share this OTP with anyone.\n\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)

	if err != nil {
		return err
	}

	fmt.Println("OTP email sent successfully")
	return nil
}
