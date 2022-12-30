package calendar

import (
	"go_cal/authentication"
	"go_cal/data"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"strconv"
	"time"
)

func UpdateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	feParams := frontendHandling.FrontendView{}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, r.Host+"/updateCalendar"))
		return
	}
	feParams, err = frontendHandling.GetFrontendParameters(r)
	if err != nil {
		cookieValue := frontendHandling.ChangeFeCookie(frontendHandling.FrontendView{})
		http.SetCookie(w, &http.Cookie{
			Name:  "fe_parameter",
			Value: cookieValue,
		})
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/updateCalendar"))
		return
	}
	if r.Method == http.MethodPost {
		switch {
		case r.Form.Has("next"):
			feParams.NextMonth()
		case r.Form.Has("prev"):
			feParams.PrevMonth()
		case r.Form.Has("today"):
			feParams.CurrentMonth()
		case r.Form.Has("choose"):
			year, err := strconv.Atoi(r.Form.Get("chooseYear"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/updateCalendar"))
				return
			}
			month, err := strconv.Atoi(r.Form.Get("chooseMonth"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/updateCalendar"))
				return
			}
			feParams.ChooseMonth(year, time.Month(month))
		}
	}
	cookieValue := frontendHandling.ChangeFeCookie(feParams)
	user := authentication.GetUserBySessionToken(r)

	http.SetCookie(w, &http.Cookie{
		Name:  "fe_parameter",
		Value: cookieValue,
	})
	templates.TempInit.Execute(w, struct {
		*frontendHandling.FrontendView
		data.User
	}{&feParams,
		*user})
	return
}
