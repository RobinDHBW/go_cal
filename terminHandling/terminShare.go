package terminHandling

import (
	"go_cal/authentication"
	"go_cal/data"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/templates"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func TerminShareHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, r.Host+"/shareTermin"))
		return
	}
	user, err := authentication.GetUserBySessionToken(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		// Fehlermeldung für Nutzer anzeigen
		templates.TempError.Execute(w, error2.CreateError(error2.Authentification, r.Host+"/"))
		return
	}
	//feParams, err := frontendHandling.GetFrontendParameters(r)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/shareTermin"))
	//	return
	//}

	switch {
	case r.Form.Has("shareTermin"):
		templates.TempShareTermin.Execute(w, user.SharedAppointments)
	case r.Form.Has("shareCreate"):
		templates.TempCreateShareTermin.Execute(w, nil)
	case r.Form.Has("terminlistShareBack"):
		templates.TempShareTermin.Execute(w, user.SharedAppointments)
	case r.Form.Has("terminShareCreateSubmit"):
		title := r.PostFormValue("title")
		if !validateInput(title) {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/shareTermin"))
			return
		}
		err := createSharedTermin(r, user, title)
		if err == (error2.DisplayedError{}) {
			templates.TempEditShareTermin.Execute(w, user.SharedAppointments[title])
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}
	case r.Form.Has("editShareTermin"):
		value := r.PostFormValue("editShareTermin")
		templates.TempEditShareTermin.Execute(w, user.SharedAppointments[value])
	case r.Form.Has("editShareTerminSubmit"):
		title := r.PostFormValue("editShareTerminSubmit")
		err := createSharedTermin(r, user, title)
		if err == (error2.DisplayedError{}) {
			templates.TempEditShareTermin.Execute(w, user.SharedAppointments[title])
		} else {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, err)
			return
		}
	case r.Form.Has("inviteUserSubmit"):
		username := r.PostFormValue("username")
		title := r.PostFormValue("inviteUserSubmit")
		if !validateInput(username) {
			w.WriteHeader(http.StatusBadRequest)
			templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, r.Host+"/shareTermin"))
			return
		}
		url := CreateURL(username, title, user.UserName)
		dataModel.Dm.AddTokenToSharedAppointment(user.Id, title, url)
		templates.TempEditShareTermin.Execute(w, user.SharedAppointments[title])
	default:
		templates.TempShareTermin.Execute(w, user.SharedAppointments)
	}
}

func createSharedTermin(r *http.Request, user *data.User, title string) error2.DisplayedError {
	begin, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateBegin"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, r.Host+"/shareTermin")
	}
	end, err := time.Parse("2006-01-02T15:04", r.Form.Get("dateEnd"))
	if err != nil {
		return error2.CreateError(error2.InvalidInput, r.Host+"/shareTermin")
	}
	if end.Before(begin) {
		return error2.CreateError(error2.EndBeforeBegin, r.Host+"/shareTermin")
	}

	repeat := GetRepeatingMode(r.Form.Get("chooseRepeat"))

	ap := data.NewAppointment(title, "", begin, end, user.Id, repeat > 0, repeat, true)
	dataModel.Dm.AddSharedAppointment(user.Id, title, ap)
	return error2.DisplayedError{}
}

func CreateURL(username, title, invitor string) string {
	token := createToken(20)
	params := url.Values{}
	params.Add("username", username)
	params.Add("termin", title)
	params.Add("token", token)
	params.Add("invitor", invitor)
	baseUrl, _ := url.Parse("http://localhost:8080/terminVoting")
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String()
}

func validateInput(text string) (successful bool) {
	// wenn Feld leer
	if len(text) == 0 {
		return false
	}
	// wenn unerlaubte Zeichen verwendet werden
	const validCharacters string = "^[a-zA-Z0-9_]*$"
	matchUsername, _ := regexp.MatchString(validCharacters, text)
	if !matchUsername {
		return false
	}
	return true
}

func createToken(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
