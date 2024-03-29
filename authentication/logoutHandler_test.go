// Matrikelnummern:
// 9495107, 4706893, 9608900

package authentication

import (
	"github.com/stretchr/testify/assert"
	"go_cal/dataModel"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogoutHandlerSuccessful(t *testing.T) {
	setup()
	defer after()
	// create user
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create session
	sessionToken, _ := CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/logout", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(LogoutHandler).ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", locationHeader.Path)

	assert.Equal(t, "session_token", response.Result().Cookies()[0].Name)
	assert.Equal(t, "", response.Result().Cookies()[0].Value)
	assert.LessOrEqual(t, response.Result().Cookies()[0].Expires, time.Now())
}

func TestLogoutHandlerNoCookie(t *testing.T) {
	setup()
	defer after()
	// create user
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create session
	CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/logout", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(LogoutHandler).ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", locationHeader.Path)

	assert.Empty(t, response.Result().Cookies())
}
