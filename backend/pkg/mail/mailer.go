package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/matcornic/hermes"
)

type Mailer struct {
	hermes *hermes.Hermes
}

func NewMailer() *Mailer {
	return &Mailer{
		hermes: &hermes.Hermes{
			Theme: new(hermes.Default),
			Product: hermes.Product{
				Name:      "Booky",
				Copyright: "Copyright © 2025 Booky. All rights reserved.",
			},
		},
	}
}

func (m *Mailer) otpEmailTemplate(otp string) (string, error) {
	e := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hey ",
			Name:     "there",
			Intros: []string{
				"We received a request to verify your email address.",
				"Use the verification code below to continue:",
				"",
				fmt.Sprintf("🔐 %s", otp),
				"",
				"This code expires in 10 minutes.",
			},
			Outros: []string{
				"If you didn't request this code, you can safely ignore this email.",
				"If you need any help, just reply to this email—we're happy to help.",
			},
		},
	}

	html, err := m.hermes.GenerateHTML(e)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML : %w", err)
	}

	return html, nil
}

func (m *Mailer) SendOTP(
	ctx context.Context,
	to, otp string,
) error {
	html, err := m.otpEmailTemplate(otp)
	if err != nil {
		return err
	}

	email := email.NewEmail()
	email.From = "hedrsag@gmail.com"
	email.To = []string{to}
	email.Subject = "Your OTP Code"
	email.HTML = []byte(html)
	email.Headers["Content-Type"] = []string{"text/html"}
	err = email.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth(
			"",
			"hedrsag@gmail.com",
			"tfhubyzpnviopmun",
			"smtp.gmail.com",
		),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
