package authentication

import (
	"github.com/stretchr/testify/assert"
	"go_cal/calendarView"
	"go_cal/dataModel"
	"go_cal/templates"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func after() {
	os.RemoveAll("../data/test/")
	os.MkdirAll("../data/test/", 777)
}

func TestSuccessfulAuthentification(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("test", "test123", 1)
	assert.Nil(t, err)
	assert.True(t, AuthenticateUser("test", "test123"))
}

func TestUnsuccessfulAuthentification(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("test", "test123", 1)
	assert.Nil(t, err)
	// wrong username
	assert.False(t, AuthenticateUser("testUser", "test123"))
	// wrong password
	assert.False(t, AuthenticateUser("test", "test"))
}

func TestCheckCookieSuccessful(t *testing.T) {
	InitServer()
	username := "test"
	// prepare session
	sessionToken, _ := createSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.True(t, checkCookie(request))
}

func TestCheckCookieUnsuccessfulWrongCookieName(t *testing.T) {
	InitServer()
	username := "test"
	// prepare session
	sessionToken, _ := createSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "wrong_session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, checkCookie(request))
}

func TestCheckCookieUnsuccessfulWrongSessionToken(t *testing.T) {
	InitServer()
	username := "test"
	// prepare session
	createSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: "cookie"})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, checkCookie(request))
}

func TestCheckCookieUnsuccessfulSessionExpired(t *testing.T) {
	InitServer()
	username := "test"
	// prepare session
	sessionToken, _ := createExpiredSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, checkCookie(request))
}

func createExpiredSession(username string) (sessionToken string, expires time.Time) {
	// Anwortchannel erstellen
	replyChannel := make(chan *session)
	// Sessiontoken generieren
	sessionToken = createUUID(25)
	// Session läuft nach x Minuten ab
	expires = time.Now().Add(-1 * time.Minute)
	// Session anhand des Sessiontokens speichern
	Serv.Cmds <- Command{ty: write, sessionToken: sessionToken, session: &session{uname: username, expires: expires}, replyChannel: replyChannel}
	// session aus Antwortchannel lesen
	session := <-replyChannel
	return sessionToken, session.expires
}

func TestCreateSession(t *testing.T) {
	InitServer()
	sessionToken, expires := createSession("testUser")
	assert.Equal(t, "testUser", GetUsernameBySessionToken(sessionToken))
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
	templates.Init()
	InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form := url.Values{}
	form.Add("uname", "testUser")
	form.Add("passwd", "test")
	form.Add("login", "")
	request.PostForm = form
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "cookie123",
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/updateCalendar", locationHeader.Path)
	cookies := response.Result().Cookies()[0]
	assert.Equal(t, "testUser", GetUsernameBySessionToken(cookies.Value))
	assert.Equal(t, "session_token", cookies.Name)
}

func TestLoginHandlerWithValidCookie(t *testing.T) {
	templates.Init()
	InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	// create User
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	sessionToken, _ := createSession("testUser")
	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/updateCalendar", locationHeader.Path)
}

func TestRegisterHandler(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	templates.Init()

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/register", nil)
	form := url.Values{}
	form.Add("uname", "testUser")
	form.Add("passwd", "test123")
	form.Add("register", "")
	request.PostForm = form

	response := httptest.NewRecorder()
	http.HandlerFunc(RegisterHandler).ServeHTTP(response, request)

	assert.Equal(t, "testUser", dataModel.Dm.GetUserByName("testUser").UserName)
	assert.True(t, dataModel.Dm.ComparePW("test123", dataModel.Dm.GetUserByName("testUser").Password))
}

func TestValidateInput(t *testing.T) {
	assert.True(t, validateInput("test_1", "test123"))
	assert.False(t, validateInput("", ""))
	assert.False(t, validateInput("test?", "test123"))
}

func TestWrapperValidCookie(t *testing.T) {
	InitServer()
	// create Session
	sessionToken, expires := createSession("testUser")
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/updateCalendar", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	time.Sleep(3 + time.Second)
	Wrapper(calendarView.UpdateCalendarHandler).ServeHTTP(response, request)
	cookies := response.Result().Cookies()[0]
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "testUser", GetUsernameBySessionToken(cookies.Value))
	assert.Equal(t, "session_token", cookies.Name)
	assert.Less(t, expires, cookies.Expires)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Calendar")
}

func TestWrapperInvalidCookie(t *testing.T) {
	InitServer()
	templates.Init()
	// create Session
	createSession("testUser")
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/updateCalendar", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "wrong_value",
	})
	response := httptest.NewRecorder()
	Wrapper(calendarView.UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Error")
}
