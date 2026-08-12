package handlers

import (
	"net/http"
	"strings"
	"time"
	"unicode"

	"Myjob/database"
	"Myjob/models"
	"Myjob/utils"

	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {

	token := strings.TrimSpace(
		r.URL.Query().Get("token"),
	)

	if token == "" {
		utils.Render(w, "reset-password.html", models.PageData{
			Title: "Reset Password",
			Error: "Invalid or missing reset token.",
		})
		return
	}

	// Find user using reset token
	user, err := database.GetUserByResetToken(token)

	if err != nil {
		http.Redirect(
			w,
			r,
			"/forgot-password",
			http.StatusSeeOther,
		)
		return
	}

	// Check if reset expiration exists
	if !user.ResetExpires.Valid ||
		user.ResetExpires.String == "" {

		utils.Render(w, "reset-password.html", models.PageData{
			Title: "Reset Password",
			Error: "Invalid or expired reset token.",
		})
		return
	}

	// Parse expiration time
	expires, err := time.Parse(
		"2006-01-02 15:04:05",
		user.ResetExpires.String,
	)

	if err != nil {
		utils.Render(w, "reset-password.html", models.PageData{
			Title: "Reset Password",
			Error: "Invalid reset token expiration.",
		})
		return
	}

	// Check if token has expired
	if time.Now().After(expires) {

		http.Redirect(
			w,
			r,
			"/forgot-password",
			http.StatusSeeOther,
		)

		return
	}

	// Show reset page
	if r.Method == "GET" {

		utils.Render(w, "reset-password.html", models.PageData{
			Title:      "Reset Password",
			ResetToken: token,
		})

		return
	}

	// Get passwords
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Validate password
	if password == "" {

		hasNumber := false
		hasUppercase := false
		hasSpecial := false

		for _, char := range password {

			if unicode.IsDigit(char) {
				hasNumber = true
			}

			if unicode.IsUpper(char) {
				hasUppercase = true
			}

			if unicode.IsPunct(char) || unicode.IsSymbol(char) {
				hasSpecial = true
			}
		}

		if len(password) < 8 {

			utils.Render(w, "reset-password.html", models.PageData{
				Title: "Reset Password",
				Error: "Password must be at least 8 characters long.",
			})

			return
		}

		if !hasNumber {

			utils.Render(w, "reset-password.html", models.PageData{
				Title: "Reset Password",
				Error: "Password must contain at least one number.",
			})

			return
		}

		if !hasUppercase {

			utils.Render(w, "reset-password.html", models.PageData{
				Title: "Reset Password",
				Error: "Password must contain at least one uppercase letter.",
			})

			return
		}

		if !hasSpecial {

			utils.Render(w, "reset-password.html", models.PageData{
				Title: "Reset Password",
				Error: "Password must contain at least one special character.",
			})

			return
		}

		utils.Render(w, "reset-password.html", models.PageData{
			Title: "Reset Password",
			Error: "Please enter a new password.",
		})

		return
	}

	if password != confirmPassword {

		utils.Render(w, "reset-password.html", models.PageData{
			Title: "Reset Password",
			Error: "Passwords do not match.",
		})

		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(
			w,
			"Unable to secure password.",
			http.StatusInternalServerError,
		)
		return
	}

	// Update password
	err = database.UpdatePassword(
		user.ID,
		string(hashedPassword),
	)

	if err != nil {
		http.Error(
			w,
			"Unable to reset password.",
			http.StatusInternalServerError,
		)
		return
	}
	err = database.ClearResetToken(user.ID)

	if err != nil {
		http.Error(
			w,
			"Password was changed, but the reset token could not be cleared.",
			http.StatusInternalServerError,
		)
		return
	}

	// Password successfully changed
	http.Redirect(
		w,
		r,
		"/login",
		http.StatusSeeOther,
	)
}
