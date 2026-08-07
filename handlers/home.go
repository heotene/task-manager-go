package handlers

import (
	"net/http"

	"Myjob/models"
	"Myjob/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {

	utils.Render(w, "home.html", models.PageData{
		Title: "Home",
	})
}
