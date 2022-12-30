package frontendHandling

import (
	"encoding/json"
	"go_cal/data"
	"net/http"
	"strings"
	"time"
)

type FrontendView struct {
	Month         time.Month
	Year          int
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
	cal.Month = time.Now().Month()
	cal.Year = time.Now().Year()
}

func (cal *FrontendView) ChooseMonth(year int, month time.Month) {
	cal.Month = month
	cal.Year = year
}

func (cal *FrontendView) GetCurrentDate() time.Time {
	return time.Now()
}

func (cal *FrontendView) GetAppointmentsForMonth(user data.User) []int {
	tl := user.Appointments
	appointmentsPerDay := make([]int, 32)
	for i := range tl {
		if tl[i].DateTimeStart.Year() == cal.Year && tl[i].DateTimeStart.Month() == cal.Month {
			appointmentsPerDay[tl[i].DateTimeStart.Day()]++
		}
		if tl[i].Timeseries.Repeat != false {
			start := tl[i].DateTimeStart
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

func GetFrontendParameters(r *http.Request) (FrontendView, error) {
	cookie, _ := r.Cookie("fe_parameter")
	fv := cookie.Value
	fv = strings.ReplaceAll(fv, "'", "\"")
	var feParams FrontendView
	err := json.Unmarshal([]byte(fv), &feParams)
	if err != nil {
		return FrontendView{}, err
	}
	return feParams, nil
}

func ChangeFeCookie(view FrontendView) string {
	if view == (FrontendView{}) {
		view = FrontendView{
			Month:         time.Now().Month(),
			Year:          time.Now().Year(),
			TerminPerSite: 7,
			TerminSite:    1,
			MinDate:       time.Now(),
		}
	}
	fvToJSON, _ := json.Marshal(view)
	fvToString := string(fvToJSON)
	fvToString = strings.ReplaceAll(fvToString, "\"", "'")
	return fvToString
}
