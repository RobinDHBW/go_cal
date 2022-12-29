package calendarAppointments

import (
	"go_cal/terminHandling"
	"time"
)

func GetAppointmentsForMonth(month time.Month, year int) []int {
	tl := terminHandling.TView.TList
	appointmentsPerDay := make([]int, 31, 31)
	for i := range (tl).Termine {
		if (tl).Termine[i].Begin.Year() == year && (tl).Termine[i].Begin.Month() == month {
			appointmentsPerDay[(tl).Termine[i].Begin.Day()-1]++
		}
		if (tl).Termine[i].Repeating != terminHandling.None {
			start := (tl).Termine[i].Begin
			// testweise nur für weekly
			for start.Before(time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)) {
				start = start.AddDate(0, 0, 7)
				if start.Year() == year && start.Month() == month {
					appointmentsPerDay[start.Day()-1]++
				}
			}
		}
	}
	return appointmentsPerDay
}
