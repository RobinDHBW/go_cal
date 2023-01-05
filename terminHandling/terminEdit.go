// Matrikelnummern:
// 9495107, 4706893, 9608900

package terminHandling

import (
	"go_cal/authentication"
	"go_cal/data"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"strconv"
	"time"
)

// TerminEditHandler
// handle inputs to edit or delete appointments
func TerminEditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/listTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentication, "/"))
		return
	}
	feParams, err := frontendHandling.GetFrontendParameters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listTermin"))
		return
	}
	switch {
	// execute template for appointment editing with appointment index given by button value
	case r.Form.Has("editTermin"):
		index, err := strconv.Atoi(r.Form.Get("editTermin"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/listTermin"))
			return
		}
		u := *user
		editIndex := GetTerminFromEditIndex(u, *feParams, index)
		templates.TempTerminEdit.Execute(w, user.Appointments[editIndex])

	// execute func to edit appointment based on user inputs
	case r.Form.Has("editTerminSubmit"):
		id, err1 := strconv.Atoi(r.FormValue("editTerminSubmit"))
		if err1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		err := EditTerminFromInput(r, true, user, id)
		errEmpty := error2.DisplayedError{}
		if err == errEmpty {
			http.Redirect(w, r, "/listTermin", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}

	// execute func to delete appointment given by index
	case r.Form.Has("deleteTerminSubmit"):
		id, err := strconv.Atoi(r.FormValue("deleteTerminSubmit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
			return
		}
		dataModel.Dm.DeleteAppointment(id, user.Id)
		http.Redirect(w, r, "/listTermin", http.StatusFound)
	default:
		http.Redirect(w, r, "/listTermin", http.StatusFound)
	}
}

// GetTerminFromEditIndex
// Calculates id of appointment that is selected to edit
// Returns this id, -1 if appointment not found
func GetTerminFromEditIndex(user data.User, fv frontendHandling.FrontendView, index int) int {
	t := fv.GetTerminList(user.Appointments)[index]
	for i := range user.Appointments {
		if user.Appointments[i].Id == t.Id {
			//currentTerminIndex = i
			return i
		}
	}
	return -1
}

// GetRepeatingMode
// calculates for given repating mode interval of appointment
func GetRepeatingMode(mode string) int {
	switch mode {
	case "none":
		return 0
	case "day":
		return 1
	case "week":
		return 7
	case "month":
		return 30
	case "year":
		return 365
	default:
		return 0
	}
}

// EditTerminFromInput
// fetch + validate frontend inputs
// edits/creates appointment with inputs, if user inputs were correct
func EditTerminFromInput(r *http.Request, edit bool, user *data.User, id int) error2.DisplayedError {
	begin, err := time.Parse("2006-01-02T15:04", r.PostFormValue("dateBegin"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, "/listTermin")
	}
	end, err := time.Parse("2006-01-02T15:04", r.PostFormValue("dateEnd"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, "/listTermin")
	}
	if end.Before(begin) {
		return error2.CreateError(error2.EndBeforeBegin, "/listTermin")
	}

	repeat := GetRepeatingMode(r.PostFormValue("chooseRepeat"))
	title := r.PostFormValue("title")
	if len(title) == 0 {
		return error2.CreateError(error2.TitleIsEmpty, "/listTermin")
	}
	content := r.PostFormValue("content")
	if edit {
		app := user.Appointments[id]
		app.Title = title
		app.Description = content
		app.DateTimeStart = begin
		app.DateTimeEnd = end
		app.Timeseries.Intervall = repeat
		app.Timeseries.Repeat = repeat > 0
		dataModel.Dm.EditAppointment(user.Id, &app)
	} else {
		//appointment := data.NewAppointment(title, content, begin, end, user.Id, repeat > 0, repeat, false)
		dataModel.Dm.AddAppointment(user.Id, title, content, "here", begin, end, repeat > 0, repeat, false)
	}
	return error2.DisplayedError{}
}
