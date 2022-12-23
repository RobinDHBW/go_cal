package terminHandling

import (
	"go_cal/calendarView"
	"go_cal/templates"
	"net/http"
	"time"
)

type Termin struct {
	Title     string
	Content   string
	Begin     time.Time
	End       time.Time
	Repeating RepeatingMode
}

type RepeatingMode string

const (
	None  RepeatingMode = "none"
	Day   RepeatingMode = "day"
	Week  RepeatingMode = "week"
	Month RepeatingMode = "month"
	Year  RepeatingMode = "year"
)

type TerminList struct {
	Termine []Termin
}

type TerminView struct {
	TList         TerminList
	TerminPerSite int
	MinDate       time.Time
}

var TView = TerminView{
	TerminList{},
	7,
	time.Now(),
}

func TerminHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch {
	case r.Form.Has("calendarBack"):
		templates.TempInit.Execute(w, calendarView.Cal)
	case r.Form.Has("dateChoose"):
		templates.TempTerminList.Execute(w, TView)
	case r.Form.Has("siteChoose"):
		templates.TempTerminList.Execute(w, TView)
	case r.Form.Has("numberPerSite"):
		templates.TempTerminList.Execute(w, TView)
	}
	templates.TempTerminList.Execute(w, TView)

}

func (tl TerminList) CreateTermin(title string, content string, begin time.Time, end time.Time) {
	termin := Termin{
		Title:     title,
		Content:   content,
		Begin:     begin,
		End:       end,
		Repeating: None,
	}
	TView.TList.Termine = append(TView.TList.Termine, termin)
}
