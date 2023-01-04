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
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/createTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, "/"))
		return
	}
	feParams, err := frontendHandling.GetFrontendParameters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/createTermin"))
		return
	}
	switch {
	case r.Form.Has("createTermin"):
		templates.TempCreateTermin.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
	case r.Form.Has("createTerminSubmit"):
		err := EditTerminFromInput(r, false, user, -1)
		errEmpty := error2.DisplayedError{}
		if err == errEmpty {
			http.Redirect(w, r, "/listTermin", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}

	default:
		templates.TempTerminList.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
	}
}
