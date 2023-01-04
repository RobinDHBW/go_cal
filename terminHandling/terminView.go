package terminHandling

import (
	"go_cal/authentication"
	"go_cal/data"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"strconv"
	"time"
)

func TerminHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, "/"))
		return
	}

	feParams, err := frontendHandling.GetFrontendParameters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
		return
	}
	switch {
	case r.Form.Has("calendarBack"):
		templates.TempInit.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
	case r.Form.Has("terminlistBack"):
		templates.TempTerminList.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
	case r.Form.Has("submitTermin"):
		input, err := strconv.Atoi(r.Form.Get("numberPerSite"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listTermin"))
			return
		}
		feParams.TerminPerSite = input

		input, err = strconv.Atoi(r.Form.Get("siteChoose"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listTermin"))
			return
		}
		feParams.TerminSite = input

		inputDate, err := time.Parse("2006-01-02", r.Form.Get("dateChoose"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listTermin"))
			return
		}
		feParams.MinDate = inputDate

		cookieValue, err := frontendHandling.GetFeCookieString(*feParams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/"))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "fe_parameter",
			Value: cookieValue,
		})
		templates.TempTerminList.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
	case r.Form.Has("searchTerminSubmit"):
		searchString := r.Form.Get("terminSearch")

		_, apps := dataModel.Dm.GetAppointmentsBySearchString(user.Id, searchString)
		templates.TempSearchTermin.Execute(w, apps)
	default:
		templates.TempTerminList.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
	}
}
