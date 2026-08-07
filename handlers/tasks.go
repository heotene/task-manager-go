package handlers

import (
	"net/http"

	"Myjob/database"
	"Myjob/middleware"
	"Myjob/models"
	"Myjob/utils"
	"strconv"
)

func Tasks(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	search := r.URL.Query().Get("search")

	var tasks []models.Task

	var err error

	if search != "" {

		tasks, err = database.SearchTasks(userID, search)

	} else {

		tasks, err = database.GetTasksByUser(userID)

	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	utils.Render(w, "tasks.html", models.PageData{
		Tasks:  tasks,
		Search: search,
	})

}

func NewTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		task := models.Task{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Priority:    r.FormValue("priority"),
			Category:    r.FormValue("category"),
			DueDate:     r.FormValue("due_date"),
			UserID:      userID,
		}

		err := database.CreateTask(task)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	utils.Render(w, "new-task.html", models.PageData{
		Title: "New Task",
	})
}
func CompleteTask(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	err := database.CompleteTask(id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	err := database.DeleteTask(id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
func EditTask(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method == "POST" {

		task := models.Task{

			ID: id,

			Title: r.FormValue("title"),

			Description: r.FormValue("description"),

			Priority: r.FormValue("priority"),

			Category: r.FormValue("category"),

			DueDate: r.FormValue("due_date"),
		}

		err := database.UpdateTask(task)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	task, err := database.GetTaskByID(id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	utils.Render(w, "edit-task.html", models.PageData{
		Title: "Edit Task",
		Task:  task,
	})
}
