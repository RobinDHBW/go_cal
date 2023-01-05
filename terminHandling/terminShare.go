// Matrikelnummern:
// 9495107, 4706893, 9608900

package terminHandling

import (
	"go_cal/authentication"
	"go_cal/data"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/templates"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TerminShareHandler handles requests due to Terminfindung
func TerminShareHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/shareTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentication, "/"))
		return
	}
	// Behandlung der verschiedenen Button-Events
	switch {
	// Terminfindung erstellen
	case r.Form.Has("shareCreate"):
		templates.TempCreateShareTermin.Execute(w, nil)
	// Eingaben zur Terminfindungserstellung bestätigen
	case r.Form.Has("terminShareCreateSubmit"):
		title := r.PostFormValue("title")
		if !validateInput(title) {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listShareTermin"))
			return
		}
		err := createSharedTermin(r, user, title)
		if err == (error2.DisplayedError{}) {
			http.Redirect(w, r, "/listShareTermin", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}
	// Terminfindung bearbeiten
	case r.Form.Has("editShareTermin"):
		value := r.PostFormValue("editShareTermin")
		templates.TempEditShareTermin.Execute(w, user.SharedAppointments[value])
	// Eingaben zur Terminfindungsbearbeitung bestätigen
	case r.Form.Has("editShareTerminSubmit"):
		title := r.PostFormValue("editShareTerminSubmit")
		err := createSharedTermin(r, user, title)
		if err == (error2.DisplayedError{}) {
			http.Redirect(w, r, "/listShareTermin", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}
	// User zu Terminfindung einladen
	case r.Form.Has("inviteUserSubmit"):
		username := r.PostFormValue("username")
		title := r.PostFormValue("inviteUserSubmit")
		if !validateInput(username) {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listShareTermin"))
			return
		}
		// URL zum Teilen erstellen
		url := dataModel.CreateURL(username, title, user.UserName)
		// eingeladenen User zur Terminfindung hinzufügen
		err := dataModel.Dm.AddTokenToSharedAppointment(user.Id, title, url, username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.DuplicateUserName, "/listShareTermin"))
			return
		}
		http.Redirect(w, r, "/listShareTermin", http.StatusFound)
	// Terminvorschlag übernehmen
	case r.Form.Has("acceptTermin"):
		// value an Trennzeichen (|) aufteilen
		parts := strings.Split(r.PostFormValue("acceptTermin"), "|")
		// Format des übergebenen values nicht valide
		if len(parts) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listShareTermin"))
			return
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listShareTermin"))
			return
		}
		// Abstimmungsergebnis wird als Beschreibung in neuen Termin übernommen
		if len(user.SharedAppointments[parts[1]]) > 0 && len(user.SharedAppointments[parts[1]]) > id {
			app := user.SharedAppointments[parts[1]][id]
			description := "Abstimmungsergebnis: \n"
			for _, val := range user.SharedAppointments[parts[1]] {
				description = description + "Zeitraum: " + val.DateTimeStart.Format("02.01.2006 15:04") + " bis " + val.DateTimeEnd.Format("02.01.2006 15:04") + " | Wiederholung: " + val.GetDescriptionFromInterval() + "\n Abgestimmt: "
				for i, voted := range val.Share.Voting {
					if voted {
						description = description + val.Share.GetUsernameFromUrl(val.Share.Tokens[i]) + ", "
					}
				}
				description = description + "\n\n"
			}
			// ausgewählten Termin übernehmen
			dataModel.Dm.AddAppointment(user.Id, app.Title, description, app.Location, app.DateTimeStart, app.DateTimeEnd, app.Timeseries.Repeat, app.Timeseries.Intervall, true)
			// Terminfindung löschen
			dataModel.Dm.DeleteSharedAppointment(app.Title, user.Id)
			http.Redirect(w, r, "/listTermin", http.StatusFound)
		}
	default:
		http.Redirect(w, r, "/listShareTermin", http.StatusFound)
	}
}

// createSharedTermin creates a shared termin for a user with the same parameters as a default appointment
func createSharedTermin(r *http.Request, user *data.User, title string) error2.DisplayedError {
	begin, err := time.Parse("2006-01-02T15:04", r.PostFormValue("dateBegin"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, "/listShareTermin")
	}
	end, err := time.Parse("2006-01-02T15:04", r.PostFormValue("dateEnd"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, "/listSareTermin")
	}
	if end.Before(begin) {
		return error2.CreateError(error2.EndBeforeBegin, "/listShareTermin")
	}
	repeat := GetRepeatingMode(r.PostFormValue("chooseRepeat"))
	dataModel.Dm.AddSharedAppointment(user.Id, title, "here", begin, end, repeat > 0, repeat, true)
	return error2.DisplayedError{}
}

// validateInput checks whether a given string is not empty and only contains alphanumerical characters (+underscore)
func validateInput(text string) (successful bool) {
	// wenn Feld leer
	if len(text) == 0 {
		return false
	}
	const validCharacters string = "^[a-zA-Z0-9_]*$"
	matchUsername, _ := regexp.MatchString(validCharacters, text)
	// wenn unerlaubte Zeichen verwendet werden
	if !matchUsername {
		return false
	}
	return true
}
