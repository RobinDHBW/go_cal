package authentication

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogoutHandler(t *testing.T) {
	deleteAllUsers()
	deleteAllSessions()

	// create user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	users["testUser"] = hashedPassword
	// create session
	sessionToken, _ := createSession("testUser")

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/logout", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(LogoutHandler).ServeHTTP(response, request)

	assert.Empty(t, sessions)

	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", locationHeader.Path)

	assert.Equal(t, "session_token", response.Result().Cookies()[0].Name)
	assert.Equal(t, "", response.Result().Cookies()[0].Value)
	assert.LessOrEqual(t, response.Result().Cookies()[0].Expires, time.Now())
}
