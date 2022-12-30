package terminHandling

import (
	"go_cal/authentication"
	"go_cal/data"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
)

func TerminCreateHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, r.Host+"/createTermin"))
		return
	}
	feParams, err := frontendHandling.GetFrontendParameters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/createTermin"))
		return
	}
	user := authentication.GetUserBySessionToken(r)
	//appointments := user.Appointments
	switch {
	case r.Form.Has("createTermin"):
		templates.TempCreateTermin.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	case r.Form.Has("createTerminSubmit"):
		if EditTerminFromInput(w, r, false, user) {
			templates.TempTerminList.Execute(w, struct {
				frontendHandling.FrontendView
				data.User
			}{feParams,
				*user})
		}

	default:
		templates.TempTerminList.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	}
}
