package notifier

import (
	"booky-backend/internal/shared/html"
	"context"
	"encoding/json"
)

type EmailHandler struct {
	renderer *html.Renderer
	mailer   Mailer
}

func NewEmailHandler(
	renderer *html.Renderer,
	mailer Mailer,
) *EmailHandler {
	return &EmailHandler{
		renderer: renderer,
		mailer:   mailer,
	}
}

func (h *EmailHandler) SendEmailOTP(ctx context.Context, msg *Message) error {
	var payload OTPPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}

	html, err := h.renderer.Render("email-otp", map[string]any{
		"Code": payload.Code,
	})
	if err != nil {
		return err
	}

	return h.mailer.SendHTML([]string{payload.Email}, "OTP Code", html)
}
