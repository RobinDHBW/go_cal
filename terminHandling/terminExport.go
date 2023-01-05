// Matrikelnummern:
// 9495107, 4706893, 9608900

package terminHandling

import (
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/export"
	"go_cal/templates"
	"net/http"
	"strconv"
)

// ICalHandler
// Handles the Download of an ICal for an user authenticated through BasicAuth
func ICalHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/listShareTermin"))
		return
	}
	uName, _, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentication, "/"))
		return
	}

	user := dataModel.Dm.GetUserByName(uName)

	//http://host/getIcal --> no queries needed
	ical := []byte(export.NewICal(&user.Appointments).ToString())

	fileName := "exportUser_" + strconv.Itoa(user.Id) + ".ics"

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "Text/Calendar")
	w.Write(ical)
}
