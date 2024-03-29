// Matrikelnummern:
// 9495107, 4706893, 9608900

package terminHandling

import (
	"go_cal/authentication"
	error2 "go_cal/error"
	"go_cal/templates"
	"net/http"
)

// TerminShareListHandler handles display of Terminfindungen
func TerminShareListHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentication, "/"))
		return
	}
	templates.TempShareTermin.Execute(w, user.SharedAppointments)
}
