package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	"go_cal/templates"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTerminCreateHandler_InvalidRequest(t *testing.T) {
	templates.Init()
	authentication.InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/editTermin", nil)
	form := url.Values{}
	request.PostForm = form
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "cookie123",
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)

	sessionToken, _ := authentication.CreateSession("testUser")
	// TODO: http und localhost
	request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/editTermin", nil)
	form = url.Values{}
	request.PostForm = form
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}
