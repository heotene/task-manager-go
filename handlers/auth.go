package handlers

import (
	"net/http"

	"Myjob/database"
	"Myjob/middleware"
	"Myjob/models"
	"Myjob/utils"

	"golang.org/x/crypto/bcrypt"
)

// When form is submitted
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		user := models.User{

			FullName: name,

			Email: email,

			Password: string(hashedPassword),
		}

		err = database.CreateUser(user)

		if err != nil {

			http.Error(w, err.Error(), 500)

			return
		}

		http.Redirect(
			w,
			r,
			"/login",
			http.StatusSeeOther,
		)

		return
	}

	utils.Render(w, "register.html", models.PageData{
		Title: "Register",
	})
}
func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Find user by email
		user, err := database.GetUserByEmail(email)

		if err != nil {
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		// Compare entered password with hashed password
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(password),
		)

		if err != nil {
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		// Create login session
		middleware.CreateSession(
			w,
			r,
			user.ID,
		)

		// Go to dashboard
		http.Redirect(
			w,
			r,
			"/dashboard",
			http.StatusSeeOther,
		)

		return
	}

	// Show login page
	utils.Render(w, "login.html", models.PageData{
		Title: "Login",
	})
}
func Dashboard(w http.ResponseWriter, r *http.Request) {

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

	total, completed, pending, err := database.GetTaskStats(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Render(w, "dashboard.html", models.PageData{
		Title:          "Dashboard",
		Name:           user.FullName,
		User:           user,
		TotalTasks:     total,
		CompletedTasks: completed,
		PendingTasks:   pending,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	middleware.DestroySession(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
