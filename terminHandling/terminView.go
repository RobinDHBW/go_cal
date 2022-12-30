package terminHandling

import (
	"go_cal/authentication"
	"go_cal/data"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func TerminHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, r.Host+"/"))
		return
	}

	feParams, _ := frontendHandling.GetFrontendParameters(r)
	switch {
	case r.Form.Has("calendarBack"):
		templates.TempInit.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	case r.Form.Has("terminlistBack"):
		templates.TempTerminList.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	case r.Form.Has("submitTermin"):
		input, err := strconv.Atoi(r.Form.Get("numberPerSite"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		feParams.TerminPerSite = input

		input, err = strconv.Atoi(r.Form.Get("siteChoose"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		feParams.TerminSite = input

		inputDate, err := time.Parse("2006-01-02", r.Form.Get("dateChoose"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		feParams.MinDate = inputDate

		cookieValue, err := frontendHandling.GetFeCookieString(feParams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/"))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "fe_parameter",
			Value: cookieValue,
		})
		//TView.GetTerminList()
		templates.TempTerminList.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	default:
		templates.TempTerminList.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	}
}

// GetTerminList
// calculates list of appointments that are later than selected date
// in case of repeating appointments, first appearance of appointment after selected date is chosen
// returns slice containing list of appointments
func GetTerminList(appointments map[int]data.Appointment, fv frontendHandling.FrontendView) []data.Appointment {
	// Create Slice to sort by date
	appSlice := make([]data.Appointment, 0, len(appointments))
	for _, i := range appointments {
		appSlice = append(appSlice, i)
	}
	sort.SliceStable(appSlice, func(i, j int) bool {
		return appSlice[i].DateTimeStart.Before(appSlice[j].DateTimeStart)
	})

	datefilteredTL := make([]data.Appointment, 0, 1)
	for i := range appSlice {
		if fv.MinDate.Before(appSlice[i].DateTimeStart) || fv.MinDate.Equal(appSlice[i].DateTimeStart) {
			datefilteredTL = append(datefilteredTL, appSlice[i])
		} else if appSlice[i].Timeseries.Repeat {
			t := GetFirstTerminOfRepeatingInDate(appSlice[i], fv)
			datefilteredTL = append(datefilteredTL, t)
		}
	}

	sort.SliceStable(datefilteredTL, func(i, j int) bool {
		return datefilteredTL[i].DateTimeStart.Before(datefilteredTL[j].DateTimeStart)
	})

	if fv.TerminPerSite*(fv.TerminSite-1) > len(datefilteredTL) {
		return nil
	}
	if fv.TerminSite*fv.TerminPerSite > len(datefilteredTL) {
		return datefilteredTL[fv.TerminPerSite*(fv.TerminSite-1):]
	}
	return datefilteredTL[fv.TerminPerSite*(fv.TerminSite-1) : fv.TerminSite*fv.TerminPerSite]
}

// GetFirstTerminOfRepeatingInDate
// calculates for given repeating appointment first appearance after selected date from FrontendView
// returns new appointment with start and end time after selected date
func GetFirstTerminOfRepeatingInDate(app data.Appointment, view frontendHandling.FrontendView) data.Appointment {
	switch app.Timeseries.Intervall {
	case 1:
		for app.DateTimeStart.Before(view.MinDate) {
			app.DateTimeStart = app.DateTimeStart.AddDate(0, 0, 1)
			app.DateTimeEnd = app.DateTimeEnd.AddDate(0, 0, 1)
		}
	case 7:
		for app.DateTimeStart.Before(view.MinDate) {
			app.DateTimeStart = app.DateTimeStart.AddDate(0, 0, 7)
			app.DateTimeEnd = app.DateTimeEnd.AddDate(0, 0, 7)
		}
	case 30:
		for app.DateTimeStart.Before(view.MinDate) {
			app.DateTimeStart = app.DateTimeStart.AddDate(0, 1, 0)
			app.DateTimeEnd = app.DateTimeEnd.AddDate(0, 1, 0)
		}
	case 365:
		for app.DateTimeStart.Before(view.MinDate) {
			app.DateTimeStart = app.DateTimeStart.AddDate(1, 0, 0)
			app.DateTimeEnd = app.DateTimeEnd.AddDate(1, 0, 0)
		}
	}
	return app
}
