// Matrikelnummern:
// 9495107, 4706893, 9608900

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

// Appointments wird für Template-Ausführung benötigt
type Appointments []data.Appointment

// TerminVotingHandler handles inputs of termin voting
// after the first GET request, the query parameters are stored in a button value in order to be transferred in a POST request
func TerminVotingHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/"))
		return
	}
	if r.Method == http.MethodPost && r.PostForm.Has("submitVoting") {
		// button value enthält zuvor gespeicherte Query parameter, getrennt durch |
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
		// herausfinden welche Checkboxen/Termine angeklickt bzw. zugesagt wurden
		for k := range r.PostForm {
			index, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			keys = append(keys, index)
		}
		user := dataModel.Dm.GetUserByName(invitor)
		// Voting ergebnis speichern
		err := dataModel.Dm.SetVotingForToken(user, keys, title, token, username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/"))
			return
		} else {
			templates.TempTerminVotingSuccess.Execute(w, nil)
		}
	} else {
		// Query Parameter auslesen
		title := r.URL.Query().Get("termin")
		invitor := r.URL.Query().Get("invitor")
		token := r.URL.Query().Get("token")
		username := r.URL.Query().Get("username")
		user := dataModel.Dm.GetUserByName(invitor)
		// Überprüfung, ob überhaupt gevotet werden darf
		if dataModel.IsVotingAllowed(title, token, user, username) {
			// Query parameter in button value schreiben, sodass sie bei einem POST ausgelesen werden können
			value := title + "|" + invitor + "|" + token + "|" + username
			aps := Appointments(dataModel.Dm.GetUserByName(invitor).SharedAppointments[title])
			templates.TempTerminVoting.Execute(w, struct {
				Appointments
				Value string
			}{aps,
				value})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidUrl, "/"))
			return
		}
	}
}
