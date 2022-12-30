package frontendHandling

import (
	"encoding/json"
	"errors"
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

// GetDaysOfMonth
// days calculation based on https://brandur.org/fragments/go-days-in-month
// calculates number of days for a month and returns slice containing the days
func (cal *FrontendView) GetDaysOfMonth() []int {
	days := time.Date(cal.Year, cal.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	dayRange := make([]int, days)
	for i := range dayRange {
		dayRange[i] = i + 1
	}
	return dayRange
}

// GetDaysBeforeMonthBegin
// calculates number of days in week before month starts
// returns empty slice with length = number of days
func (cal *FrontendView) GetDaysBeforeMonthBegin() []int {
	weekday := time.Date(cal.Year, cal.Month, 1, 0, 0, 0, 0, time.UTC).Weekday()
	if weekday == 0 {
		return make([]int, 6)
	} else {
		return make([]int, weekday-1)
	}
}

// NextMonth
// adds 1 month to given FrontendView
func (cal *FrontendView) NextMonth() {
	if cal.Month == 12 {
		cal.Month = 1
		cal.Year++
	} else {
		cal.Month++
	}
}

// PrevMonth
// subtracts one month from given FrontendView
func (cal *FrontendView) PrevMonth() {
	if cal.Month == 1 {
		cal.Month = 12
		cal.Year--
	} else {
		cal.Month--
	}
}

// CurrentMonth
// set month and year from FrontendView to current date
func (cal *FrontendView) CurrentMonth() {
	cal.Month = time.Now().Month()
	cal.Year = time.Now().Year()
}

// ChooseMonth
// set month and year to custom, error if date is not valid
func (cal *FrontendView) ChooseMonth(year int, month time.Month) error {
	if year < 0 || month < 1 || month > 12 {
		return errors.New("given date not valid")
	}
	cal.Month = month
	cal.Year = year
	return nil
}

// GetCurrentDate
// returns current datetime
func (cal *FrontendView) GetCurrentDate() time.Time {
	return time.Now()
}

// GetAppointmentsForMonth
// searches for appointments for given user and month
// returns slice with number of appointments for each day
func (cal *FrontendView) GetAppointmentsForMonth(user data.User) []int {
	tl := user.Appointments
	appointmentsPerDay := make([]int, 32) //max. index = 31 (number of days)
	for i := range tl {
		if tl[i].DateTimeStart.Year() == cal.Year && tl[i].DateTimeStart.Month() == cal.Month {
			appointmentsPerDay[tl[i].DateTimeStart.Day()]++
		}
		if tl[i].Timeseries.Repeat {
			start := tl[i].DateTimeStart
			switch tl[i].Timeseries.Intervall {
			case 1:
				for start.Before(time.Date(cal.Year, cal.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(0, 0, 1)
					if start.Year() == cal.Year && start.Month() == cal.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			case 7:
				for start.Before(time.Date(cal.Year, cal.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(0, 0, 7)
					if start.Year() == cal.Year && start.Month() == cal.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			case 30:
				for start.Before(time.Date(cal.Year, cal.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(0, 1, 0)
					if start.Year() == cal.Year && start.Month() == cal.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			case 365:
				for start.Before(time.Date(cal.Year, cal.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(1, 0, 0)
					if start.Year() == cal.Year && start.Month() == cal.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			}
		}
	}
	return appointmentsPerDay
}

// GetFrontendParameters
// returns FrontendView struct out of Cookie
func GetFrontendParameters(r *http.Request) (FrontendView, error) {
	cookie, err := r.Cookie("fe_parameter")
	if err != nil {
		return FrontendView{}, err
	}
	fv := cookie.Value
	fv = strings.ReplaceAll(fv, "'", "\"")
	var feParams FrontendView
	err = json.Unmarshal([]byte(fv), &feParams)
	if err != nil {
		return FrontendView{}, err
	}
	return feParams, nil
}

// GetFeCookieString
// returns json-string of FrontendView to store in Cookie
// creates new FrontendView if parameter-struct is empty
func GetFeCookieString(view FrontendView) (string, error) {
	if view == (FrontendView{}) {
		view = FrontendView{
			Month:         time.Now().Month(),
			Year:          time.Now().Year(),
			TerminPerSite: 7,
			TerminSite:    1,
			MinDate:       time.Now(),
		}
	}
	fvToJSON, err := json.Marshal(view)
	if err != nil {
		return "", err
	}
	fvToString := string(fvToJSON)
	fvToString = strings.ReplaceAll(fvToString, "\"", "'")
	return fvToString, nil
}
