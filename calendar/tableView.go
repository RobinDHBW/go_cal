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

// UpdateCalendarHandler
// handle inputs of the calendar view
// change FrontendView-Cookie based on inputs, changing calendar-View
func UpdateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	feParams := &frontendHandling.FrontendView{}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/updateCalendar"))
		return
	}
	// Get Cookie or set default-Cookie if not existing
	feParams, err = frontendHandling.GetFrontendParameters(r)
	if err != nil {
		cookieValue, err := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/"))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "fe_parameter",
			Value: cookieValue,
		})
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/updateCalendar"))
		return
	}
	if r.Method == http.MethodPost {
		switch {
		// Set calendar to next month
		case r.Form.Has("next"):
			feParams.NextMonth()
		// Set calendar to previous month
		case r.Form.Has("prev"):
			feParams.PrevMonth()
		// Set calendar to current month
		case r.Form.Has("today"):
			feParams.CurrentMonth()
		// Get values from inputs
		case r.Form.Has("choose"):
			year, err := strconv.Atoi(r.Form.Get("chooseYear"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/updateCalendar"))
				return
			}
			month, err := strconv.Atoi(r.Form.Get("chooseMonth"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/updateCalendar"))
				return
			}
			err = feParams.ChooseMonth(year, time.Month(month))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/updateCalendar"))
				return
			}
		}
	}
	cookieValue, err := frontendHandling.GetFeCookieString(*feParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, "/"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "fe_parameter",
		Value: cookieValue,
	})
	templates.TempInit.Execute(w, struct {
		*frontendHandling.FrontendView
		*data.User
	}{feParams,
		user})
	return
}
