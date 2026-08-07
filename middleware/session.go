package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore(
	[]byte("my-secret-key"),
)

func CreateSession(w http.ResponseWriter, r *http.Request, userID int) {

	session, _ := Store.Get(r, "session")

	session.Values["user_id"] = userID

	session.Save(r, w)
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
