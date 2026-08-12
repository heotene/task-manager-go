package handlers

import (
	"net/http"

	"Myjob/database"
	"Myjob/middleware"
	"Myjob/models"
	"Myjob/utils"
)

func DeleteAccount(w http.ResponseWriter, r *http.Request) {

	// Make sure the user is logged in.
	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Show confirmation page.
	if r.Method == "GET" {

		utils.Render(w, "delete-account.html", models.PageData{
			Title: "Delete Account",
		})

		return
	}

	// Only accept POST for actual deletion.
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Require the confirmation phrase.
	confirmation := r.FormValue("confirmation")

	if confirmation != "DELETE MY ACCOUNT" {

		utils.Render(w, "delete-account.html", models.PageData{
			Title: "Delete Account",
			Error: "Please type DELETE MY ACCOUNT exactly to confirm.",
		})

		return
	}

	// Delete account and tasks.
	err := database.DeleteAccount(userID)

	if err != nil {
		http.Error(
			w,
			"Unable to delete account: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// Destroy the login session.
	middleware.DestroySession(w, r)

	// Send the user back home.
	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}
