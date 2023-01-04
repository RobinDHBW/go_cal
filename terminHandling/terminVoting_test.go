package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestTerminVotingHandlerUnsuccessfulNoForm(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodPost, "/terminVoting", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Default2))
}

func TestTerminVotingHandlerGetQueryVotingNotAllowed(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodGet, "/terminVoting?invitor=anna&termin=test&token=GzybyHccrBcnAmlmOAOt&username=peter", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.InvalidUrl))
}

func TestTerminVotingHandlerGetQueryVotingAllowed(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("anna", "test", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = dataModel.Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("peter", "Terminfindung1", "anna"), "peter")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	expectedValueInButton := "Terminfindung1|anna|" + token + "|peter"
	request, _ := http.NewRequest(http.MethodGet, "/terminVoting?invitor=anna&termin=Terminfindung1&token="+token+"&username=peter", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Titel: Terminfindung1")
	assert.Contains(t, string(body), "Terminvorschläge")
	assert.Contains(t, string(body), "03.01.2023 22:00")
	assert.Contains(t, string(body), expectedValueInButton)
}

func TestTerminVotingHandlerPostWrongValue(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("anna", "test", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = dataModel.Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("peter", "Terminfindung1", "anna"), "peter")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	request, _ := http.NewRequest(http.MethodPost, "/terminVoting?invitor=anna&termin=Terminfindung1&token="+token+"&username=peter", nil)
	form := url.Values{}
	// wrong button value, username peter is missing
	form.Add("submitVoting", "Terminfindung1|anna|"+token)
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.InvalidInput))
}

func TestTerminVotingHandlerPostNoUser(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("anna", "test", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = dataModel.Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("peter", "Terminfindung1", "anna"), "peter")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	request, _ := http.NewRequest(http.MethodPost, "/terminVoting?invitor=anna&termin=Terminfindung1&token="+token+"&username=peter", nil)
	form := url.Values{}
	// wrong button value, token is corrupted
	form.Add("submitVoting", "Terminfindung1|anna|"+token+"s|peter")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Default2))
}

func TestTerminVotingHandlerPostCorrectValue(t *testing.T) {
	setup()
	defer after()
	user, err := dataModel.Dm.AddUser("anna", "test", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	dataModel.Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = dataModel.Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("peter", "Terminfindung1", "anna"), "peter")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	// initial ist der Terminvorschlag abgesagt
	assert.False(t, user.SharedAppointments["Terminfindung1"][0].Share.Voting[0])
	request, _ := http.NewRequest(http.MethodPost, "/terminVoting?invitor=anna&termin=Terminfindung1&token="+token+"&username=peter", nil)
	form := url.Values{}
	form.Add("submitVoting", "Terminfindung1|anna|"+token+"|peter")
	// 1. Termin zusagen
	form.Add("0", "on")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Erfolgreich abgestimmt.")
	// Termin wurde zugesagt
	assert.True(t, user.SharedAppointments["Terminfindung1"][0].Share.Voting[0])
}
