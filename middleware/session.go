package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore(
	[]byte("my-secret-key"),
)

func CreateSession(w http.ResponseWriter, r *http.Request, userID int) {

	session, err := Store.Get(r, "session")

	if err != nil {
		http.Error(w, "Unable to create session", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = userID

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, "Unable to save session", http.StatusInternalServerError)
		return
	}
}

func GetUserID(r *http.Request) int {

	session, _ := Store.Get(r, "session")

	userID, ok := session.Values["user_id"].(int)

	if !ok {
		return 0
	}

	return userID
}

func DestroySession(w http.ResponseWriter, r *http.Request) {

	session, _ := Store.Get(r, "session")

	session.Options.MaxAge = -1

	session.Save(r, w)
}

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userID := GetUserID(r)

		if userID <= 0 {
			http.Redirect(
				w,
				r,
				"/login",
				http.StatusSeeOther,
			)
			return
		}

		next(w, r)
	}
}
