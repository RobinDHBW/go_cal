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

var currentTerminIndex int = -1

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
		//fmt.Println(index)
		u := *user
		GetTerminFromEditIndex(u, feParams, index)
		templates.TempTerminEdit.Execute(w, user.Appointments[currentTerminIndex])

	case r.Form.Has("editTerminSubmit"):
		if EditTerminFromInput(w, r, true, user) {
			templates.TempTerminList.Execute(w, struct {
				frontendHandling.FrontendView
				data.User
			}{feParams,
				*user})

		}

	case r.Form.Has("deleteTerminSubmit"):
		DeleteTermin(user)
		templates.TempTerminList.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})

	default:
		templates.TempTerminList.Execute(w, struct {
			frontendHandling.FrontendView
			data.User
		}{feParams,
			*user})
	}
}

func GetTerminFromEditIndex(user data.User, fv frontendHandling.FrontendView, index int) int {
	t := GetTerminList(user, fv)[index]
	for i := range user.Appointments {
		if user.Appointments[i].Id == t.Id {
			currentTerminIndex = i
			return i
		}
	}
	return -1
}

func editTermin(app *data.Appointment, title string, content string, begin time.Time, end time.Time, repeat int) {
	app.Title = title
	app.Description = content
	app.DateTimeStart = begin
	app.DateTimeEnd = end
	app.Timeseries.Intervall = repeat
	app.Timeseries.Repeat = repeat > 0
}

func DeleteTermin(user *data.User) {
	dataModel.Dm.DeleteAppointment(user.Appointments[currentTerminIndex].Id, user.Id)
	//(*tl).Termine[currentTerminIndex] = (*tl).Termine[len((*tl).Termine)-1]
	//currentTerminIndex = -1
	//(*tl).Termine = (*tl).Termine[:len((*tl).Termine)-1]
}

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

func EditTerminFromInput(w http.ResponseWriter, r *http.Request, edit bool, user *data.User) bool {
	begin, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateBegin"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
		return false
	}
	end, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateEnd"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/listTermin"))
		return false
	}
	if end.Before(begin) {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.EndBeforeBegin, r.Host+"/listTermin"))
		return false
	}

	repeat := GetRepeatingMode(r.Form.Get("chooseRepeat"))
	title := r.Form.Get("title")
	if len(title) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		templates.TempError.Execute(w, error2.CreateError(error2.TitleIsEmpty, r.Host+"/listTermin"))
		return false
	}
	content := r.Form.Get("content")
	if edit {
		app := user.Appointments[currentTerminIndex]
		//TODO checken ob das so geht
		editTermin(&app, title, content, begin, end, repeat)
		dataModel.Dm.EditAppointment(user.Id, app)
	} else {
		appointment := data.NewAppointment(title, content, begin, end, user.Id, repeat > 0, repeat, false, "")
		dataModel.Dm.AddAppointment(user.Id, appointment)
	}
	return true
}
