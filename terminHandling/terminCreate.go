package terminHandling

import (
	error2 "go_cal/error"
	"go_cal/templates"
	"net/http"
)

func TerminCreateHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		error2.CreateError(error2.Default2, "/createTermin", w, http.StatusInternalServerError)
		return
	}

	switch {
	case r.Form.Has("createTermin"):
		templates.TempCreateTermin.Execute(w, TView)
	case r.Form.Has("createTerminSubmit"):
		if TView.TList.EditTerminFromInput(w, r, false) {
			templates.TempTerminList.Execute(w, TView)
		}

	default:
		templates.TempTerminList.Execute(w, TView)
	}
}
