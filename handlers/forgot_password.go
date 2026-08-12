package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"Myjob/database"
	"Myjob/models"
	"Myjob/utils"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		login := strings.TrimSpace(r.FormValue("login"))
		token := strings.TrimSpace(r.FormValue("token"))

		// STEP 1: Generate reset code
		if token == "" {

			if login == "" {
				utils.Render(w, "forgot-password.html", models.PageData{
					Title: "Forgot Password",
					Error: "Please enter your email or phone number.",
				})
				return
			}

			// Determine verification method
			verificationMethod := "phone"

			if strings.Contains(login, "@") {
				verificationMethod = "email"
			}

			// Try email first
			user, err := database.GetUserByEmail(login)

			// If email was not found, try phone
			if err != nil {
				user, err = database.GetUserByPhone(login)
			}

			if err != nil {
				utils.Render(w, "forgot-password.html", models.PageData{
					Title: "Forgot Password",
					Error: "No account was found with that email or phone number.",
				})
				return
			}

			// Generate 6-digit reset code
			tokenBytes := make([]byte, 4)

			_, err = rand.Read(tokenBytes)

			if err != nil {
				http.Error(
					w,
					"Unable to generate reset code.",
					http.StatusInternalServerError,
				)
				return
			}

			token := fmt.Sprintf(
				"%06d",
				uint32(tokenBytes[0])<<24|
					uint32(tokenBytes[1])<<16|
					uint32(tokenBytes[2])<<8|
					uint32(tokenBytes[3]),
			)

			token = token[len(token)-6:]

			// Token expires in 15 minutes
			expires := time.Now().
				Add(15 * time.Minute).
				Format("2006-01-02 15:04:05")

			// Save reset code
			err = database.UpdateResetToken(
				user.ID,
				token,
				expires,
			)

			if err != nil {
				http.Error(
					w,
					"Unable to create password reset request.",
					http.StatusInternalServerError,
				)
				return
			}

			// Send reset code by email
			if verificationMethod == "email" {

				err = utils.SendResetCode(
					user.Email,
					token,
				)

				if err != nil {

					log.Println("EMAIL ERROR:", err)

					utils.Render(w, "forgot-password.html", models.PageData{
						Title: "Forgot Password",
						Error: "Unable to send reset code. Please try again.",
					})

					return
				}
			}

			// Show token input page
			utils.Render(w, "forgot-password.html", models.PageData{
				Title:              "Forgot Password",
				Success:            "Reset token sent. Enter the token below.",
				Search:             login,
				VerificationMethod: verificationMethod,
			})

			return
		}

		// STEP 2: Token was entered
		http.Redirect(
			w,
			r,
			"/reset-password?token="+url.QueryEscape(token),
			http.StatusSeeOther,
		)

		return
	}

	// GET request
	utils.Render(w, "forgot-password.html", models.PageData{
		Title: "Forgot Password",
	})
}
