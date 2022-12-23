package calendarView

import "time"

type Calendar struct {
	Month   time.Month
	Year    int
	Current time.Time
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
	} else {
		return make([]int, weekday-1)
	}
}

func (cal Calendar) NextMonth() {
	if cal.Month == 12 {
		cal.Month = 0
		cal.Year++
	} else {
		cal.Month++
	}
}

func (cal Calendar) PrevMonth() {
	if cal.Month == 0 {
		cal.Month = 12
		cal.Year--
	} else {
		cal.Month--
	}
}

func (cal Calendar) CurrentMonth() {
	cal.Month = cal.Current.Month()
	cal.Year = cal.Current.Year()
}

func (cal Calendar) ChooseMonth(year int, month time.Month) {
	cal.Month = month
	cal.Year = year
}
