package terminHandling

import (
	error2 "go_cal/error"
	"go_cal/templates"
	"net/http"
	"sort"
	"strconv"
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
	None  RepeatingMode = "keine"
	Day   RepeatingMode = "täglich"
	Week  RepeatingMode = "wöchentlich"
	Month RepeatingMode = "monatlich"
	Year  RepeatingMode = "jährlich"
)

type TerminList struct {
	Termine []Termin
}

type TerminView struct {
	TList         TerminList
	TerminPerSite int
	TerminSite    int
	MinDate       time.Time
}

var Tlist = TerminList{}

var TView = TerminView{
	Tlist,
	7,
	1,
	time.Now(),
}

func TerminHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
		return
	}

	switch {
	case r.Form.Has("calendarBack"):
		//ToDo wieder einkommentieren (aber dann error wegen cycle import)
		//templates.TempInit.Execute(w, calendarView.Cal)
	case r.Form.Has("terminlistBack"):
		templates.TempTerminList.Execute(w, TView)
	case r.Form.Has("submitTermin"):
		input, err := strconv.Atoi(r.Form.Get("numberPerSite"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		TView.TerminPerSite = input

		input, err = strconv.Atoi(r.Form.Get("siteChoose"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		TView.TerminSite = input

		inputDate, err := time.Parse("2006-01-02", r.Form.Get("dateChoose"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		TView.MinDate = inputDate

		//TView.GetTerminList()
		templates.TempTerminList.Execute(w, TView)
	default:
		templates.TempTerminList.Execute(w, TView)
	}

}

func (tl *TerminList) CreateTermin(title string, content string, begin time.Time, end time.Time, repeat RepeatingMode) {
	termin := Termin{
		Title:     title,
		Content:   content,
		Begin:     begin,
		End:       end,
		Repeating: repeat,
	}
	tl.Termine = append(tl.Termine, termin)
}

func (tv TerminView) GetTerminList() []Termin {
	sort.SliceStable(tv.TList.Termine, func(i, j int) bool {
		return tv.TList.Termine[i].Begin.Before(tv.TList.Termine[j].Begin)
	})

	datefilteredTL := make([]Termin, 0, 1)
	for i := range tv.TList.Termine {
		if tv.MinDate.Before(tv.TList.Termine[i].Begin) || tv.MinDate.Equal(tv.TList.Termine[i].Begin) {
			datefilteredTL = append(datefilteredTL, tv.TList.Termine[i])
		} else if tv.TList.Termine[i].Repeating != None {
			t := tv.TList.Termine[i].GetFirstTerminOfRepeatingInDate()
			datefilteredTL = append(datefilteredTL, t)
		}
	}

	sort.SliceStable(datefilteredTL, func(i, j int) bool {
		return datefilteredTL[i].Begin.Before(datefilteredTL[j].Begin)
	})

	if tv.TerminPerSite*(tv.TerminSite-1) > len(datefilteredTL) {
		return nil
	}
	if tv.TerminSite*tv.TerminPerSite > len(datefilteredTL) {
		return datefilteredTL[tv.TerminPerSite*(tv.TerminSite-1):]
	}
	return datefilteredTL[tv.TerminPerSite*(tv.TerminSite-1) : tv.TerminSite*tv.TerminPerSite]
}

func (t Termin) GetFirstTerminOfRepeatingInDate() Termin {
	switch t.Repeating {
	case Week:
		for t.Begin.Before(TView.MinDate) {
			t.Begin = t.Begin.AddDate(0, 0, 7)
			t.End = t.End.AddDate(0, 0, 7)
		}
		return t
	case Day:
		for t.Begin.Before(TView.MinDate) {
			t.Begin = t.Begin.AddDate(0, 0, 1)
			t.End = t.End.AddDate(0, 0, 1)
		}
		return t
	case Month:
		for t.Begin.Before(TView.MinDate) {
			t.Begin = t.Begin.AddDate(0, 1, 0)
			t.End = t.End.AddDate(0, 1, 0)
		}
	case Year:
		for t.Begin.Before(TView.MinDate) {
			t.Begin = t.Begin.AddDate(1, 0, 0)
			t.End = t.End.AddDate(1, 0, 0)
		}
		return t
	}
	return t
}
