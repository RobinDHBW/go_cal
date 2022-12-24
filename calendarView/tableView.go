package calendarView

import (
	"go_cal/authentication"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Calendar struct {
	Month   time.Month
	Year    int
	Current time.Time
}

// has to be removed
var Cal = Calendar{
	Month:   time.Now().Month(),
	Year:    time.Now().Year(),
	Current: time.Now(),
}

// https://brandur.org/fragments/go-days-in-month
func (cal Calendar) GetDaysOfMonth() []int {
	days := time.Date(cal.Year, cal.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	dayRange := make([]int, days)
	for i := range dayRange {
		dayRange[i] = i + 1
	}
	return dayRange
}

func (cal Calendar) GetDaysBeforeMonthBegin() []int {
	weekday := time.Date(cal.Year, cal.Month, 1, 0, 0, 0, 0, time.UTC).Weekday()
	if weekday == 6 {
		return make([]int, 5)
	} else if weekday == 0 {
		return make([]int, 0)
	} else {
		return make([]int, weekday-1)
	}
}

func (cal *Calendar) NextMonth() {
	if cal.Month == 12 {
		cal.Month = 1
		cal.Year++
	} else {
		cal.Month++
	}
}

func (cal *Calendar) PrevMonth() {
	if cal.Month == 1 {
		cal.Month = 12
		cal.Year--
	} else {
		cal.Month--
	}
}

func (cal *Calendar) CurrentMonth() {
	cal.Month = cal.Current.Month()
	cal.Year = cal.Current.Year()
}

func (cal *Calendar) ChooseMonth(year int, month time.Month) {
	cal.Month = month
	cal.Year = year
}

func UpdateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	isCookieValid := authentication.CheckCookie(r)
	if isCookieValid {
		if r.Method == http.MethodPost {
			r.ParseForm()
			switch {
			case r.Form.Has("next"):
				Cal.NextMonth()
			case r.Form.Has("prev"):
				Cal.PrevMonth()
			case r.Form.Has("today"):
				Cal.CurrentMonth()
			case r.Form.Has("choose"):
				year, _ := strconv.Atoi(r.Form.Get("chooseYear"))
				month, _ := strconv.Atoi(r.Form.Get("chooseMonth"))
				Cal.ChooseMonth(year, time.Month(month))
			}
		}
	} else {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/error?type=authentification&link="+url.QueryEscape("/"), http.StatusUnauthorized)
		return
	}

	var tempInit = template.Must(template.ParseFiles("./templates/test.tmpl.html"))
	tempInit.Execute(w, Cal)
}
