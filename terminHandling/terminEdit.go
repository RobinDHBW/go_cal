package terminHandling

import (
	"fmt"
	"go_cal/templates"
	"net/http"
	"strconv"
)

func TerminEditHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	switch {
	case r.Form.Has("editTermin"):
		index, _ := strconv.Atoi(r.Form.Get("editTermin"))
		TView.TList.GetTerminFromEditIndex(index)
		templates.TempTerminEdit.Execute(w, TView)
	default:
		templates.TempTerminList.Execute(w, TView)
	}
}

func (tl *TerminList) GetTerminFromEditIndex(index int) {
	t := TView.GetTerminList()[index]
	fmt.Println(t)
}
