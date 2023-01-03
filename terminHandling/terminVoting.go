package terminHandling

import (
	"go_cal/data"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/templates"
	"net/http"
	"strconv"
	"strings"
)

type Appointments []data.Appointment

func TerminVotingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == http.MethodPost && r.PostForm.Has("submitVoting") {
		split := strings.Split(r.PostFormValue("submitVoting"), "|")
		if len(split) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/"))
			return
		}
		title := split[0]
		invitor := split[1]
		token := split[2]
		username := split[3]
		keys := make([]int, 0, len(r.PostForm))
		for k := range r.PostForm {
			index, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			keys = append(keys, index)
		}
		user := dataModel.Dm.GetUserByName(invitor)
		if user != nil {
			err := dataModel.Dm.SetVotingForToken(user, keys, title, token, username)
			if err == nil {
				templates.TempTerminVotingSuccess.Execute(w, nil)
			}
		}
	} else {
		title := r.URL.Query().Get("termin")
		invitor := r.URL.Query().Get("invitor")
		token := r.URL.Query().Get("token")
		username := r.URL.Query().Get("username")
		user := dataModel.Dm.GetUserByName(invitor)
		if user != nil {
			if dataModel.IsVotingAllowed(title, token, user, username) {
				// Query parameter in button value schreiben, sodass sie bei einem POST ausgelesen werden können
				value := title + "|" + invitor + "|" + token + "|" + username
				aps := Appointments(dataModel.Dm.GetUserByName(invitor).SharedAppointments[title])
				templates.TempTerminVoting.Execute(w, struct {
					Appointments
					Value string
				}{aps,
					value})
			}
		}
	}
}
