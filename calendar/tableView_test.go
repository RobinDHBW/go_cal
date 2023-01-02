package calendar

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	"go_cal/templates"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TODo: funktioniert noch nicht
func TestUpdateCalendarHandler(t *testing.T) {
	defer after()
	authentication.InitServer()
	templates.Init()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("Testuser", "test", 0)
	assert.Nil(t, err)
	sessionToken, _ := authentication.CreateSession("Testuser")

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/updateCalendar", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)

	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/updateCalendar", locationHeader.Path)
}

func after() {
	os.RemoveAll("../data/test/")
	os.MkdirAll("../data/test/", 777)
}
