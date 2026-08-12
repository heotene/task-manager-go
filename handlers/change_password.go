package handlers

import (
	"net/http"

	"Myjob/database"
	"Myjob/middleware"
	"Myjob/models"
	"Myjob/utils"

	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := database.GetUserByID(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {

		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		// Check current password
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(currentPassword),
		)

		if err != nil {
			renderChangePasswordError(
				w,
				"Current password is incorrect.",
			)
			return
		}

		// Check new password
		if len(newPassword) < 6 {
			renderChangePasswordError(
				w,
				"New password must be at least 6 characters.",
			)
			return
		}

		// Confirm password
		if newPassword != confirmPassword {
			renderChangePasswordError(
				w,
				"New passwords do not match.",
			)
			return
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(newPassword),
			bcrypt.DefaultCost,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Update password
		err = database.UpdateUserPassword(
			userID,
			string(hashedPassword),
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		modelsData := models.PageData{
			Title:   "Change Password",
			User:    user,
			Success: "Password changed successfully.",
		}

		utils.Render(
			w,
			"change-password.html",
			modelsData,
		)

		return
	}

	utils.Render(
		w,
		"change-password.html",
		models.PageData{
			Title: "Change Password",
			User:  user,
		},
	)
}

func renderChangePasswordError(
	w http.ResponseWriter,
	message string,
) {
	utils.Render(
		w,
		"change-password.html",
		models.PageData{
			Title: "Change Password",
			Error: message,
		},
	)
}
