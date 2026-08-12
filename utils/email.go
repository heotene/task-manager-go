package utils

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendResetCode(email string, code string) error {

	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not configured")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Task Manager <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Your Task Manager Password Reset Code",
		Html: fmt.Sprintf(`
			<h2>Password Reset</h2>

			<p>Your Task Manager password reset code is:</p>

			<h1>%s</h1>

			<p>This code will expire in 15 minutes.</p>

			<p>If you did not request a password reset, you can ignore this email.</p>
		`, code),
	}

	_, err := client.Emails.Send(params)

	return err
}

func SendVerificationCode(email string, code string) error {

	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not configured")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Task Manager <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Your Task Manager Verification Code",
		Html: fmt.Sprintf(`
			<h2>Verify Your Account</h2>

			<p>Your Task Manager verification code is:</p>

			<h1>%s</h1>

			<p>This code expires in 10 minutes.</p>

			<p>If you did not create this account, you can ignore this email.</p>
		`, code),
	}

	_, err := client.Emails.Send(params)

	return err
}
