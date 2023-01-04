package frontendHandling

import (
	"encoding/json"
	"errors"
	"go_cal/data"
	"net/http"
	"sort"
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
func (fv *FrontendView) GetDaysOfMonth() []int {
	days := time.Date(fv.Year, fv.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	dayRange := make([]int, days)
	for i := range dayRange {
		dayRange[i] = i + 1
	}
	return dayRange
}

// GetDaysBeforeMonthBegin
// calculates number of days in week before month starts
// returns empty slice with length = number of days
func (fv *FrontendView) GetDaysBeforeMonthBegin() []int {
	weekday := time.Date(fv.Year, fv.Month, 1, 0, 0, 0, 0, time.UTC).Weekday()
	if weekday == 0 {
		return make([]int, 6)
	} else {
		return make([]int, weekday-1)
	}
}

// NextMonth
// adds 1 month to given FrontendView
func (fv *FrontendView) NextMonth() {
	if fv.Month == 12 {
		fv.Month = 1
		fv.Year++
	} else {
		fv.Month++
	}
}

// PrevMonth
// subtracts one month from given FrontendView
func (fv *FrontendView) PrevMonth() {
	if fv.Month == 1 {
		fv.Month = 12
		fv.Year--
	} else {
		fv.Month--
	}
}

// CurrentMonth
// set month and year from FrontendView to current date
func (fv *FrontendView) CurrentMonth() {
	fv.Month = time.Now().Month()
	fv.Year = time.Now().Year()
}

// ChooseMonth
// set month and year to custom, error if date is not valid
func (fv *FrontendView) ChooseMonth(year int, month time.Month) error {
	if year < 0 || month < 1 || month > 12 {
		return errors.New("given date not valid")
	}
	fv.Month = month
	fv.Year = year
	return nil
}

// GetCurrentDate
// returns current datetime
func (fv *FrontendView) GetCurrentDate() time.Time {
	return time.Now()
}

// GetAppointmentsForMonth
// searches for appointments for given user and month
// returns slice with number of appointments for each day
func (fv *FrontendView) GetAppointmentsForMonth(user data.User) []int {
	tl := user.Appointments
	appointmentsPerDay := make([]int, 32) //max. index = 31 (number of days)
	for i := range tl {
		if tl[i].DateTimeStart.Year() == fv.Year && tl[i].DateTimeStart.Month() == fv.Month {
			appointmentsPerDay[tl[i].DateTimeStart.Day()]++
		}
		if tl[i].Timeseries.Repeat {
			start := tl[i].DateTimeStart
			switch tl[i].Timeseries.Intervall {
			case 1:
				for start.Before(time.Date(fv.Year, fv.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(0, 0, 1)
					if start.Year() == fv.Year && start.Month() == fv.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			case 7:
				for start.Before(time.Date(fv.Year, fv.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(0, 0, 7)
					if start.Year() == fv.Year && start.Month() == fv.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			case 30:
				for start.Before(time.Date(fv.Year, fv.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(0, 1, 0)
					if start.Year() == fv.Year && start.Month() == fv.Month {
						appointmentsPerDay[start.Day()]++
					}
				}
			case 365:
				for start.Before(time.Date(fv.Year, fv.Month+1, 1, 0, 0, 0, 0, time.Local)) {
					start = start.AddDate(1, 0, 0)
					if start.Year() == fv.Year && start.Month() == fv.Month {
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
func GetFrontendParameters(r *http.Request) (*FrontendView, error) {
	cookie, err := r.Cookie("fe_parameter")
	if err != nil {
		return &FrontendView{}, err
	}
	fv := cookie.Value
	fv = strings.ReplaceAll(fv, "'", "\"")
	var feParams FrontendView
	err = json.Unmarshal([]byte(fv), &feParams)
	if err != nil {
		return &FrontendView{}, err
	}
	return &feParams, nil
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
			MinDate:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 00, 0, 0, time.UTC),
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

// GetTerminList
// calculates list of appointments that are later than selected date
// in case of repeating appointments, first appearance of appointment after selected date is chosen
// returns slice containing list of appointments
func (fv *FrontendView) GetTerminList(appointments map[int]data.Appointment) []data.Appointment {
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
			t := GetFirstTerminOfRepeatingInDate(appSlice[i], *fv)
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
func GetFirstTerminOfRepeatingInDate(app data.Appointment, view FrontendView) data.Appointment {
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
