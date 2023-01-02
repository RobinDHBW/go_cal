package calendar

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// TODo: funktioniert noch nicht
func TestUpdateCalendarHandler(t *testing.T) {
	authentication.Serv = &authentication.Server{Cmds: authentication.StartSessionManager()}
	dataModel.InitDataModel()
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)

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
	_ = os.Remove("../files/" + strconv.FormatInt(int64(user.Id), 10) + ".json")

}
