package utils

import (
	"html/template"
	"net/http"

	"Myjob/models"
)

func Render(w http.ResponseWriter, page string, data models.PageData) {

	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/"+page,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

}
