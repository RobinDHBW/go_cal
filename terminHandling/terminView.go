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
	user := authentication.GetUserBySessionToken(r)
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

		cookieValue := frontendHandling.ChangeFeCookie(feParams)
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

//func (tl *TerminList) CreateTermin(title string, content string, begin time.Time, end time.Time, repeat RepeatingMode) {
//	termin := Termin{
//		Title:     title,
//		Content:   content,
//		Begin:     begin,
//		End:       end,
//		Repeating: repeat,
//	}
//	tl.Termine = append(tl.Termine, termin)
//}

func GetTerminList(user data.User, fv frontendHandling.FrontendView) []data.Appointment {
	sort.SliceStable(user.Appointments, func(i, j int) bool {
		return user.Appointments[i].DateTimeStart.Before(user.Appointments[j].DateTimeStart)
	})

	datefilteredTL := make([]data.Appointment, 0, 1)
	for i := range user.Appointments {
		if fv.MinDate.Before(user.Appointments[i].DateTimeStart) || fv.MinDate.Equal(user.Appointments[i].DateTimeStart) {
			datefilteredTL = append(datefilteredTL, user.Appointments[i])
		} else if user.Appointments[i].Timeseries.Repeat {
			t := GetFirstTerminOfRepeatingInDate(user.Appointments[i], fv)
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
