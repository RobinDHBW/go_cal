package authentication

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestSuccessfulAuthentification(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	users["testUser"] = hashedPassword
	assert.True(t, AuthenticateUser("testUser", []byte("test123")))
}

func TestUnsuccessfulAuthentification(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	users["testUser"] = hashedPassword
	// wrong username
	assert.False(t, AuthenticateUser("wrongTestUserName", []byte("test123")))
	// wrong password
	assert.False(t, AuthenticateUser("testUser", []byte("test")))
}

func TestDuplicateUsername(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	users["testUser"] = hashedPassword
	// duplicate username
	assert.True(t, isDuplicateUsername("testUser"))
	// no duplicate username
	assert.False(t, isDuplicateUsername("test"))
}

func TestCheckCookieSuccessful(t *testing.T) {
	username := "testUser"
	sessionToken := "cookie123"
	expires := time.Now().Add(120 * time.Second)
	// prepare session
	sessions[sessionToken] = session{
		uname:   username,
		expires: expires,
	}
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.True(t, CheckCookie(request))
}

func TestCheckCookieUnsuccessfulWrongCookieName(t *testing.T) {
	username := "testUser"
	sessionToken := "cookie123"
	expires := time.Now().Add(120 * time.Second)
	// prepare session
	sessions[sessionToken] = session{
		uname:   username,
		expires: expires,
	}
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "wrong_session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, CheckCookie(request))
}

func TestCheckCookieUnsuccessfulWrongSessionToken(t *testing.T) {
	username := "testUser"
	sessionToken := "cookie123"
	expires := time.Now().Add(120 * time.Second)
	// prepare session
	sessions[sessionToken] = session{
		uname:   username,
		expires: expires,
	}
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: "cookie"})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, CheckCookie(request))
}

func TestCheckCookieUnsuccessfulSessionExpired(t *testing.T) {
	username := "testUser"
	sessionToken := "cookie123"
	expires := time.Now().Add(-120 * time.Second)
	// prepare session
	sessions[sessionToken] = session{
		uname:   username,
		expires: expires,
	}
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: "cookie123"})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, CheckCookie(request))
}

func TestCreateSession(t *testing.T) {
	// initially no sessions
	assert.Empty(t, sessions)
	sessionToken, expires := createSession("testUser")
	assert.NotEmpty(t, sessions)
	assert.Equal(t, "testUser", sessions[sessionToken].uname)
	assert.LessOrEqual(t, expires.Sub(time.Now()).Minutes(), 2.0)
}

func TestIsExpired(t *testing.T) {
	session := session{
		uname:   "testUser",
		expires: time.Now().Add(120 * time.Second),
	}
	assert.False(t, session.isExpired())

	session.expires = time.Now().Add(-120 * time.Second)
	assert.True(t, session.isExpired())
}

func TestLoginHandlerWithoutCookie(t *testing.T) {
	deleteAllUsers()
	deleteAllSessions()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	users["testUser"] = hashedPassword

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form := url.Values{}
	form.Add("uname", "testUser")
	form.Add("passwd", "test123")
	form.Add("login", "")
	request.PostForm = form
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "cookie123",
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	location, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/updateCalendar", location.Path)
	cookies := response.Result().Cookies()[0]
	_, ok := sessions[cookies.Value]
	assert.True(t, ok)
	assert.Equal(t, "testUser", sessions[cookies.Value].uname)
	assert.Equal(t, "session_token", cookies.Name)

	assert.Equal(t, sessions[cookies.Value].expires.UTC().Round(1*time.Second), cookies.Expires.UTC().Round(1*time.Second))
}

func TestLoginHandlerWithCookie(t *testing.T) {
	username := "testUser"
	sessionToken := "cookie123"
	expires := time.Now().Add(120 * time.Second)
	// prepare session
	sessions[sessionToken] = session{
		uname:   username,
		expires: expires,
	}

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("uname=testUser&passwd=test&login=")))
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	//client := &http.Client{}
	//response, _ := client.Do(request)
}

// TODO: filepaths not working
//func TestLoadUsersFromFiles(t *testing.T) {
//	assert.NoError(t, LoadUsersFromFiles())
//	assert.NotEmpty(t, users)
//}

func TestLogoutHandler(t *testing.T) {

}

func TestRegisterHandler(t *testing.T) {

}

func deleteAllUsers() {
	for k := range users {
		delete(users, k)
	}
}

func deleteAllSessions() {
	for k := range sessions {
		delete(sessions, k)
	}
}
