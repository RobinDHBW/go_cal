package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestTerminShareHandlerUnsuccessfulNoForm(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Default2))
}

func TestTerminShareHandlerUnsuccessfulNoSessiontoken(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Authentification))
}

// shareCreate-Button
func TestTerminShareHandlerSuccessfulShareCreate(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung erstellen-Button
	form.Add("shareCreate", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Terminfindung erstellen")
	assert.Contains(t, string(body), "Erstellen")
}

// terminShareCreateSubmit-Button
func TestTerminShareHandlerSuccessfulTerminShareCreateSubmit(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung erstellen-Button
	form.Add("terminShareCreateSubmit", "")
	form.Add("title", "Terminfindung1")
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listShareTermin", locationHeader.Path)
}

// terminShareCreateSubmit-Button
func TestTerminShareHandlerUnsuccessfulTerminShareCreateSubmitEndBeforeBegin(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung erstellen-Button
	form.Add("terminShareCreateSubmit", "")
	form.Add("title", "Terminfindung1")
	// Enddatum vor Anfangsdatum --> error
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T21:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.EndBeforeBegin))
}

// terminShareCreateSubmit-Button
func TestTerminShareHandlerUnsuccessfulTerminShareCreateSubmitEmptyTitle(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung erstellen-Button
	form.Add("terminShareCreateSubmit", "")
	// leerer Titel --> error
	form.Add("title", "")
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.InvalidInput))
}

// editShareTermin-Button
func TestTerminShareHandlerSuccessfulEditShareTermin(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung bearbeiten-Button
	form.Add("editShareTermin", "Terminfindung1")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Bisherige Terminvorschläge:")
	assert.Contains(t, string(body), "Terminfindung1")
	assert.Contains(t, string(body), "Bisherige eingeladene User:")
	assert.Contains(t, string(body), "Keine User eingeladen.")
	assert.Contains(t, string(body), "Terminvorschläge hinzufügen")
	assert.Contains(t, string(body), "User einladen")
}

// editShareTerminSubmit-Button
func TestTerminShareHandlerSuccessfulEditShareTerminSubmit(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung bearbeiten-Button
	form.Add("editShareTerminSubmit", "Terminfindung1")
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listShareTermin", locationHeader.Path)
}

// editShareTerminSubmit-Button
func TestTerminShareHandlerUnsuccessfulEditShareTerminSubmitEndBeforeBegin(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// Terminfindung bearbeiten-Button
	form.Add("editShareTerminSubmit", "Terminfindung1")
	// Enddatum vor Anfangsdatum --> error
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T21:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.EndBeforeBegin))
}

// inviteUserSubmit-Button
func TestTerminShareHandlerSuccessfulInviteUserSubmit(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// initial keine User eingeladen
	assert.Equal(t, 0, len(user.SharedAppointments["Terminfindung1"][0].Share.Tokens))
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// User zu Terminfindung einladen-Button
	form.Add("inviteUserSubmit", "Terminfindung1")
	form.Add("username", "Anna")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listShareTermin", locationHeader.Path)
	// ein User wurde eingeladen
	assert.Equal(t, 1, len(user.SharedAppointments["Terminfindung1"][0].Share.Tokens))
}

// inviteUserSubmit-Button
func TestTerminShareHandlerUnsuccessfulInviteUserSubmitInvalidInput(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// User zu Terminfindung einladen-Button
	form.Add("inviteUserSubmit", "Terminfindung1")
	// Username leer --> error
	form.Add("username", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.InvalidInput))
}

// inviteUserSubmit-Button
func TestTerminShareHandlerUnsuccessfulInviteUserSubmitDuplicateUsername(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// User zu Terminfindung einladen-Button
	form.Add("inviteUserSubmit", "Terminfindung1")
	form.Add("username", "Anna")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	// erster User einladen funktioniert ohne Fehler
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listShareTermin", locationHeader.Path)

	// zweiten User mit dem gleichen Usernamen einladen --> error
	request, _ = http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form = url.Values{}
	// User zu Terminfindung einladen-Button
	form.Add("inviteUserSubmit", "Terminfindung1")
	form.Add("username", "Anna")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.DuplicateUserName))
}

// acceptTermin-Button
func TestTerminShareHandlerSuccessfulAcceptTermin(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	assert.Empty(t, user.Appointments)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	assert.Equal(t, 1, len(user.SharedAppointments))
	assert.Equal(t, 0, len(user.Appointments))
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	form.Add("acceptTermin", "0|Terminfindung1")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, 0, len(user.SharedAppointments))
	assert.Equal(t, 1, len(user.Appointments))
	// SharedTermin hatte ID 1
	assert.Equal(t, "Terminfindung1", user.Appointments[2].Title)
	assert.Contains(t, user.Appointments[2].Description, "Abgestimmt:")
	assert.Contains(t, user.Appointments[2].Description, "Abstimmungsergebnis:")
	assert.Equal(t, beginDate, user.Appointments[2].DateTimeStart)
	assert.Equal(t, endDate, user.Appointments[2].DateTimeEnd)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listTermin", locationHeader.Path)
}

// acceptTermin-Button
func TestTerminShareHandlerUnsuccessfulWrongValue(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	assert.Empty(t, user.Appointments)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// more than 2 splits
	form.Add("acceptTermin", "0|Terminfindung1|test")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.InvalidInput))

}

// acceptTermin-Button
func TestTerminShareHandlerUnsuccessfulInvalidId(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	assert.Empty(t, user.Appointments)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	// more than 2 splits
	form.Add("acceptTermin", "-1|Terminfindung1")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), string(error2.InvalidInput))
}

// default-case
func TestTerminShareHandlerDefaultCase(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	form := url.Values{}
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listShareTermin", locationHeader.Path)
}

func TestCreateSharedTerminSuccessful(t *testing.T) {
	setup()
	defer after()
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, error2.DisplayedError{}, displayedError)
	assert.Equal(t, 1, len(user.SharedAppointments["Terminfindung"]))
	assert.Equal(t, "Terminfindung", user.SharedAppointments["Terminfindung"][0].Title)
	assert.Equal(t, user.Id, user.SharedAppointments["Terminfindung"][0].Userid)
	parsed, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	assert.Equal(t, parsed, user.SharedAppointments["Terminfindung"][0].DateTimeStart)
	parsed, _ = time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	assert.Equal(t, parsed, user.SharedAppointments["Terminfindung"][0].DateTimeEnd)
}

func TestCreateSharedTerminUnsuccessfulWrongFormatBegin(t *testing.T) {
	setup()
	defer after()
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	// wrong format in dateBegin
	form.Add("dateBegin", "2023-0103T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, string(error2.InvalidInput), displayedError.Text)
}

func TestCreateSharedTerminUnsuccessfulWrongFormatEnd(t *testing.T) {
	setup()
	defer after()
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	// wrong format in dateEnd
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-0103T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, string(error2.InvalidInput), displayedError.Text)
}

func TestCreateSharedTerminUnsuccessfulEndBeforeBegin(t *testing.T) {
	setup()
	defer after()
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	// wrong format in dateEnd
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T21:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, string(error2.EndBeforeBegin), displayedError.Text)
}

func TestValidateInput(t *testing.T) {
	successful := validateInput("")
	assert.False(t, successful)
	successful = validateInput("test?")
	assert.False(t, successful)
	successful = validateInput("Test123_")
	assert.True(t, successful)
}
