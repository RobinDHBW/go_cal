package calendarView

import (
	"encoding/json"
	error2 "go_cal/error"
	"go_cal/templates"
	"go_cal/terminHandling"
	"net/http"
	"strconv"
	"time"
)

type FrontendView struct {
	Month         time.Month
	Year          int
	Current       time.Time
	TerminPerSite int
	TerminSite    int
	MinDate       time.Time
}

// https://brandur.org/fragments/go-days-in-month
func (cal *FrontendView) GetDaysOfMonth() []int {
	days := time.Date(cal.Year, cal.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	dayRange := make([]int, days)
	for i := range dayRange {
		dayRange[i] = i + 1
	}
	return dayRange
}

func (cal *FrontendView) GetDaysBeforeMonthBegin() []int {
	weekday := time.Date(cal.Year, cal.Month, 1, 0, 0, 0, 0, time.UTC).Weekday()
	if weekday == 0 {
		return make([]int, 6)
	} else {
		return make([]int, weekday-1)
	}
}

func (cal *FrontendView) NextMonth() {
	if cal.Month == 12 {
		cal.Month = 1
		cal.Year++
	} else {
		cal.Month++
	}
}

func (cal *FrontendView) PrevMonth() {
	if cal.Month == 1 {
		cal.Month = 12
		cal.Year--
	} else {
		cal.Month--
	}
}

func (cal *FrontendView) CurrentMonth() {
	cal.Month = cal.Current.Month()
	cal.Year = cal.Current.Year()
}

func (cal *FrontendView) ChooseMonth(year int, month time.Month) {
	cal.Month = month
	cal.Year = year
}

func (cal *FrontendView) GetAppointmentsForMonth() []int {
	tl := terminHandling.TView.TList
	appointmentsPerDay := make([]int, 32)
	for i := range (tl).Termine {
		if (tl).Termine[i].Begin.Year() == cal.Year && (tl).Termine[i].Begin.Month() == cal.Month {
			appointmentsPerDay[(tl).Termine[i].Begin.Day()]++
		}
		if (tl).Termine[i].Repeating != terminHandling.None {
			start := (tl).Termine[i].Begin
			// testweise nur für weekly
			for start.Before(time.Date(cal.Year, cal.Month+1, 1, 0, 0, 0, 0, time.Local)) {
				start = start.AddDate(0, 0, 7)
				if start.Year() == cal.Year && start.Month() == cal.Month {
					appointmentsPerDay[start.Day()]++
				}
			}
		}
	}
	return appointmentsPerDay
}

func GetFrontendParameters(w http.ResponseWriter, r *http.Request) FrontendView {
	cookie, _ := r.Cookie("fe_parameter")
	fv := cookie.Value
	var feParams FrontendView
	err := json.Unmarshal([]byte(fv), &feParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/updateCalendar"))
		return FrontendView{}
	}
	return feParams
}

func ChangeFeCookie(w http.ResponseWriter, view FrontendView) {
	fvToJSON, _ := json.Marshal(view)
	http.SetCookie(w, &http.Cookie{
		Name:  "fe_parameter",
		Value: string(fvToJSON),
	})
}

func UpdateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	feParams := FrontendView{}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/updateCalendar"))
			return
		}
		feParams = GetFrontendParameters(w, r)

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
		ChangeFeCookie(w, feParams)
	}
	templates.TempInit.Execute(w, &feParams)
	return
}
