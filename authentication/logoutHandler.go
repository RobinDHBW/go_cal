package authentication

import (
	"net/http"
	"time"
)

// LogoutHandler handle logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	replyChannel := make(chan *session)
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	sessionToken := cookie.Value
	// Session entfernen
	Serv.Cmds <- Command{ty: remove, sessionToken: sessionToken, replyChannel: replyChannel}
	<-replyChannel
	// Info für den Client
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "fe_parameter",
		Value: "",
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
