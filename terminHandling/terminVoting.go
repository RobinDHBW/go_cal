package terminHandling

import (
	"go_cal/dataModel"
	"go_cal/templates"
	"net/http"
)

func TerminVotingHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	title := r.URL.Query().Get("termin")
	//token := r.URL.Query().Get("token")
	invitor := r.URL.Query().Get("invitor")
	r.ParseForm()
	if r.Method == http.MethodPost && r.PostForm.Has("submitVoting") {
		// funkioniert noch nicht
		http.Redirect(w, r, CreateURL(username, title, invitor), http.StatusFound)
	}
	templates.TempTerminVoting.Execute(w, dataModel.Dm.GetUserByName(invitor).SharedAppointments[title])
}
