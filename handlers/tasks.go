package handlers

import (
	"net/http"

	"Myjob/database"
	"Myjob/middleware"
	"Myjob/models"
	"Myjob/utils"
	"strconv"
	"strings"
)

func Tasks(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	search := r.URL.Query().Get("search")
	priority := r.URL.Query().Get("priority")
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")

	tasks, err := database.SearchAndFilterTasks(
		userID,
		search,
		priority,
		category,
		status,
		sort,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Render(w, "tasks.html", models.PageData{
		Title:    "My Tasks",
		Tasks:    tasks,
		Search:   search,
		Priority: priority,
		Category: category,
		Status:   status,
		Sort:     sort,
	})
}
func NewTask(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {

		title := strings.TrimSpace(r.FormValue("title"))
		description := strings.TrimSpace(r.FormValue("description"))
		priority := strings.TrimSpace(r.FormValue("priority"))
		category := strings.TrimSpace(r.FormValue("category"))
		dueDate := strings.TrimSpace(r.FormValue("due_date"))

		// Validate title
		if title == "" {
			utils.Render(w, "new-task.html", models.PageData{
				Title: "New Task",
				Error: "Task title is required.",
			})
			return
		}

		// Validate priority
		if priority != "High" &&
			priority != "Medium" &&
			priority != "Low" {

			utils.Render(w, "new-task.html", models.PageData{
				Title: "New Task",
				Error: "Please select a valid priority.",
			})
			return
		}

		task := models.Task{
			Title:       title,
			Description: description,
			Priority:    priority,
			Category:    category,
			DueDate:     dueDate,
			UserID:      userID,
		}

		err := database.CreateTask(task)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = database.CompleteTask(id, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteTask(id, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func EditTask(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {

		title := strings.TrimSpace(r.FormValue("title"))
		description := strings.TrimSpace(r.FormValue("description"))
		priority := strings.TrimSpace(r.FormValue("priority"))
		category := strings.TrimSpace(r.FormValue("category"))
		dueDate := strings.TrimSpace(r.FormValue("due_date"))

		// Validate title
		if title == "" {
			task, _ := database.GetTaskByIDAndUser(id, userID)

			utils.Render(w, "edit-task.html", models.PageData{
				Title: "Edit Task",
				Task:  task,
				Error: "Task title is required.",
			})
			return
		}

		// Validate priority
		if priority != "High" &&
			priority != "Medium" &&
			priority != "Low" {

			task, _ := database.GetTaskByIDAndUser(id, userID)

			utils.Render(w, "edit-task.html", models.PageData{
				Title: "Edit Task",
				Task:  task,
				Error: "Please select a valid priority.",
			})
			return
		}

		task := models.Task{
			ID:          id,
			Title:       title,
			Description: description,
			Priority:    priority,
			Category:    category,
			DueDate:     dueDate,
			UserID:      userID,
		}

		err = database.UpdateTask(task)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	task, err := database.GetTaskByIDAndUser(id, userID)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	utils.Render(w, "edit-task.html", models.PageData{
		Title: "Edit Task",
		Task:  task,
	})
}
