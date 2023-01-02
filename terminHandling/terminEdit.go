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

func TerminEditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, r.Host+"/editTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, r.Host+"/"))
		return
	}
	feParams, err := frontendHandling.GetFrontendParameters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/editTermin"))
		return
	}
	switch {
	case r.Form.Has("editTermin"):
		index, err := strconv.Atoi(r.Form.Get("editTermin"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/editTermin"))
			return
		}
		u := *user
		editIndex := GetTerminFromEditIndex(u, *feParams, index)
		templates.TempTerminEdit.Execute(w, user.Appointments[editIndex])

	case r.Form.Has("editTerminSubmit"):
		id, _ := strconv.Atoi(r.FormValue("editTerminSubmit")) //ToDo testen ob das so geht
		err := EditTerminFromInput(r, true, user, id)
		errEmpty := error2.DisplayedError{}
		if err == errEmpty {
			templates.TempTerminList.Execute(w, struct {
				*frontendHandling.FrontendView
				*data.User
			}{feParams,
				user})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}

	case r.Form.Has("deleteTerminSubmit"):
		id, _ := strconv.Atoi(r.FormValue("deleteTerminSubmit")) //ToDo testen ob das so geht

		//editIndex := GetTerminFromEditIndex(u, *feParams, index)
		dataModel.Dm.DeleteAppointment(id, user.Id)

		//DeleteTermin(user)
		templates.TempTerminList.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})

	default:
		templates.TempTerminList.Execute(w, struct {
			*frontendHandling.FrontendView
			*data.User
		}{feParams,
			user})
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

//func DeleteTermin(user *data.User) {
//	dataModel.Dm.DeleteAppointment(user.Appointments[currentTerminIndex].Id, user.Id)
//	//(*tl).Termine[currentTerminIndex] = (*tl).Termine[len((*tl).Termine)-1]
//	//currentTerminIndex = -1
//	//(*tl).Termine = (*tl).Termine[:len((*tl).Termine)-1]
//}

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
	begin, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateBegin"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, r.Host+"/listTermin")
	}
	end, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateEnd"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, r.Host+"/listTermin")
	}
	if end.Before(begin) {
		return error2.CreateError(error2.EndBeforeBegin, r.Host+"/listTermin")
	}

	repeat := GetRepeatingMode(r.Form.Get("chooseRepeat"))
	title := r.Form.Get("title")
	if len(title) == 0 {
		return error2.CreateError(error2.TitleIsEmpty, r.Host+"/listTermin")
	}
	content := r.Form.Get("content")
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
