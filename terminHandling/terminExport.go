package terminHandling

import (
	"go_cal/authentication"
	error2 "go_cal/error"
	"go_cal/export"
	"go_cal/templates"
	"net/http"
	"strconv"
)

func ICalHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, r.Host+"/listShareTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, r.Host+"/"))
		return
	}

	//http://host/getIcal --> no queries needed
	ical := []byte(export.NewICal(&user.Appointments).ToString())

	fileName := "exportUser_" + strconv.Itoa(user.Id) + ".ics"

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "Text/Calendar")
	w.Write(ical)
}
