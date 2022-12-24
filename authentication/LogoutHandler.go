package authentication

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		return
	}
	sessionToken := cookie.Value

	// user session aus sessions-map entfernen
	delete(sessions, sessionToken)

	// Info für den Client
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
