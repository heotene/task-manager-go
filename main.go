package main

import (
	"log"
	"net/http"

	"Myjob/database"
	"Myjob/handlers"
	"Myjob/middleware"
)

func main() {

	// Connect to database
	database.Connect()

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.Home)

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/verify", handlers.Verify)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc(
		"/dashboard",
		middleware.RequireLogin(handlers.Dashboard),
	)

	http.HandleFunc("/tasks", handlers.Tasks)
	http.HandleFunc("/tasks/new", handlers.NewTask)
	http.HandleFunc("/tasks/edit", handlers.EditTask)
	http.HandleFunc("/tasks/complete", handlers.CompleteTask)
	http.HandleFunc("/tasks/delete", handlers.DeleteTask)

	http.HandleFunc("/profile", handlers.Profile)

	http.HandleFunc("/settings", handlers.Settings)

	http.HandleFunc("/change-password", handlers.ChangePassword)

	http.HandleFunc("/settings/delete-account", handlers.DeleteAccount)

	http.HandleFunc("/forgot-password", handlers.ForgotPassword)

	http.HandleFunc("/reset-password", handlers.ResetPassword)

	// Start server
	log.Println("Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
