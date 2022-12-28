package terminHandling

import (
	error2 "go_cal/error"
	"go_cal/templates"
	"net/http"
	"strconv"
	"time"
)

var currentTerminIndex int = -1

func TerminEditHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		error2.CreateError(error2.Default2, "/terminedit", w, http.StatusInternalServerError)
		return
	}

	switch {
	case r.Form.Has("editTermin"):
		index, err := strconv.Atoi(r.Form.Get("editTermin"))
		if err != nil {
			error2.CreateError(error2.InvalidInput, "/terminedit", w, http.StatusBadRequest)
			return
		}
		TView.TList.GetTerminFromEditIndex(index)
		templates.TempTerminEdit.Execute(w, TView.TList.Termine[currentTerminIndex])

	case r.Form.Has("editTerminSubmit"):
		TView.TList.EditTerminFromInput(w, r, true)
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

func (tl *TerminList) EditTerminFromInput(w http.ResponseWriter, r *http.Request, edit bool) {
	begin, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateBegin"))
	if err != nil {
		error2.CreateError(error2.InvalidInput, "/terminlist", w, http.StatusBadRequest)
		return
	}
	end, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateEnd"))
	if err != nil {
		error2.CreateError(error2.InvalidInput, "/terminlist", w, http.StatusBadRequest)
		return
	}
	if end.Before(begin) {
		error2.CreateError(error2.EndBeforeBegin, "/terminlist", w, http.StatusBadRequest)
		return
	}

	repeat := GetRepeatingMode(r.Form.Get("chooseRepeat"))
	title := r.Form.Get("title")
	if len(title) == 0 {
		error2.CreateError(error2.TitleIsEmpty, "/terminlist", w, http.StatusBadRequest)
		return
	}
	content := r.Form.Get("content")
	if edit {
		tl.Termine[currentTerminIndex].editTermin(title, content, begin, end, repeat)
	} else {
		tl.CreateTermin(title, content, begin, end, repeat)
	}
}
