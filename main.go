package main

import (
	"log"
	"net/http"

	"Myjob/database"
	"Myjob/handlers"
)

func main() {

	// Connect to database
	database.Connect()

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.Home)

	http.HandleFunc("/login", handlers.Login)

	http.HandleFunc("/register", handlers.Register)

	http.HandleFunc("/dashboard", handlers.Dashboard)

	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc("/tasks", handlers.Tasks)

	http.HandleFunc("/tasks/new", handlers.NewTask)

	http.HandleFunc("/tasks/complete", handlers.CompleteTask)

	http.HandleFunc("/tasks/delete", handlers.DeleteTask)

	http.HandleFunc("/tasks/edit", handlers.EditTask)

	// Start server
	log.Println("Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
