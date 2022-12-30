package authentication

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	replyChannel := make(chan *session)
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	sessionToken := cookie.Value

	Serv.Cmds <- Command{ty: remove, sessionToken: sessionToken, replyChannel: replyChannel}
	<-replyChannel
	// user session aus sessions-map entfernen
	//delete(sessions, sessionToken)

	// Info für den Client
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
