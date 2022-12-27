package terminHandling

import (
	"go_cal/templates"
	"net/http"
	"strconv"
	"time"
)

var currentTerminIndex int = -1

func TerminEditHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	switch {
	case r.Form.Has("editTermin"):
		index, _ := strconv.Atoi(r.Form.Get("editTermin"))
		TView.TList.GetTerminFromEditIndex(index)
		templates.TempTerminEdit.Execute(w, TView.TList.Termine[currentTerminIndex])
	case r.Form.Has("editTerminSubmit"):
		begin, _ := time.Parse("2006-01-02T15:04", r.Form.Get("dateBegin"))
		end, _ := time.Parse("2006-01-02T15:04", r.Form.Get("dateEnd"))
		repeat := GetRepeatingMode(r.Form.Get("chooseRepeat"))
		TView.TList.Termine[currentTerminIndex].editTermin(r.Form.Get("title"), r.Form.Get("content"), begin, end, repeat)
		templates.TempTerminList.Execute(w, TView)
	case r.Form.Has("deleteTerminSubmit"):
		TView.TList.DeleteTermin()
		templates.TempTerminList.Execute(w, TView)
	default:
		templates.TempTerminList.Execute(w, TView)
	}
}

func (tl *TerminList) GetTerminFromEditIndex(index int) {
	t := TView.GetTerminList()[index]
	for i := range (*tl).Termine {
		if t == (*tl).Termine[i] {
			currentTerminIndex = i
		}
	}
}

func (t *Termin) editTermin(title string, content string, begin time.Time, end time.Time, repeat RepeatingMode) {
	t.Title = title
	t.Content = content
	t.Begin = begin
	t.End = end
	t.Repeating = repeat
}

func (tl *TerminList) DeleteTermin() {
	(*tl).Termine[currentTerminIndex] = (*tl).Termine[len((*tl).Termine)-1]
	currentTerminIndex = -1
	(*tl).Termine = (*tl).Termine[:len((*tl).Termine)-1]
}

func GetRepeatingMode(mode string) RepeatingMode {
	switch mode {
	case "none":
		return None
	case "day":
		return Day
	case "week":
		return Week
	case "month":
		return Month
	case "year":
		return Year
	default:
		return None
	}
}
