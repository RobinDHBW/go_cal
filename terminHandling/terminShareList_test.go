package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/templates"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTerminShareListHandlerUnsuccessful(t *testing.T) {
	authentication.InitServer()
	templates.Init()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	request, _ := http.NewRequest(http.MethodPost, "/listShareTermin", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareListHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	// Fehler, da kein Cookie im Request
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Authentification))
}

func TestTerminShareListHandlerSuccessful(t *testing.T) {
	authentication.InitServer()
	templates.Init()
	defer after()
	dataModel.InitDataModel("../data/test")
	user, err := dataModel.Dm.AddUser("peter", "test", 1)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("peter")
	request, _ := http.NewRequest(http.MethodPost, "/listShareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminShareListHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	// keine Terminfindungen erstellt
	assert.Contains(t, string(body), "Keine Terminfindungen vorhanden")

	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	request, _ = http.NewRequest(http.MethodPost, "/listShareTermin", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminShareListHandler).ServeHTTP(response, request)
	body, _ = io.ReadAll(response.Result().Body)
	// Terminfindung wird angezeigt
	assert.Contains(t, string(body), "Terminfindung1")
	assert.Contains(t, string(body), "Terminfindung anzeigen")
}