package authentication

import (
	"github.com/stretchr/testify/assert"
	"go_cal/configuration"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/templates"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

const dataPath = "../data/test/authentication"

func after() {
	err := os.RemoveAll(dataPath)
	if err != nil {
		return
	}
}

func TestSuccessfulAuthentification(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("test", "test123", 1)
	assert.Nil(t, err)
	assert.True(t, AuthenticateUser("test", "test123"))
}

func TestUnsuccessfulAuthentification(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("test", "test123", 1)
	assert.Nil(t, err)
	// wrong username
	assert.False(t, AuthenticateUser("testUser", "test123"))
	// wrong password
	assert.False(t, AuthenticateUser("test", "test"))
}

func TestCheckCookieSuccessful(t *testing.T) {
	setup()
	defer after()
	username := "test"
	// prepare session
	sessionToken, _ := CreateSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.True(t, checkCookie(request))
}

func TestCheckCookieUnsuccessfulWrongCookieName(t *testing.T) {
	setup()
	defer after()
	username := "test"
	// prepare session
	sessionToken, _ := CreateSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "wrong_session_token", Value: sessionToken})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, checkCookie(request))
}

func TestCheckCookieUnsuccessfulWrongSessionToken(t *testing.T) {
	setup()
	defer after()
	username := "test"
	// prepare session
	CreateSession(username)
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "session_token", Value: "cookie"})
	// copy cookie to request
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.False(t, checkCookie(request))
}

func TestCheckCookieUnsuccessfulSessionExpired(t *testing.T) {
	setup()
	defer after()
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
	setup()
	defer after()
	sessionToken, expires := CreateSession("testUser")
	assert.Equal(t, "testUser", getUsernameBySessionToken(sessionToken))
	assert.LessOrEqual(t, expires.Sub(time.Now()).Minutes(), 10.0)
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
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
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
	assert.Equal(t, "testUser", getUsernameBySessionToken(cookies.Value))
	assert.Equal(t, "session_token", cookies.Name)
}

func TestLoginHandlerWithValidCookie(t *testing.T) {
	setup()
	defer after()
	// create User
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	sessionToken, _ := CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
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

func TestLoginHandlerNoForm(t *testing.T) {
	setup()
	defer after()
	// create User
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Default2))
}

func TestLoginHandlerInvalidInput(t *testing.T) {
	setup()
	defer after()
	// create User
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Add("uname", "")
	form.Add("passwd", "test123")
	form.Add("login", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.EmptyField))
}

func TestLoginHandlerAuthenticationFailed(t *testing.T) {
	setup()
	defer after()
	// create User
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Add("uname", "wrongUsername")
	form.Add("passwd", "test123")
	form.Add("login", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.WrongCredentials))
}

func TestLoginHandlerGetTemplate(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, string(body), "Login")
}

func TestRegisterHandler(t *testing.T) {
	setup()
	defer after()

	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
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

func TestRegisterHandlerNoForm(t *testing.T) {
	setup()
	defer after()

	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(RegisterHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.Default2))
}

func TestRegisterHandlerInvalidInput(t *testing.T) {
	setup()
	defer after()

	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
	form := url.Values{}
	form.Add("uname", "")
	form.Add("passwd", "test123")
	form.Add("register", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(RegisterHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.EmptyField))
}

func TestRegisterHandlerDuplicateUsername(t *testing.T) {
	setup()
	defer after()

	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
	form := url.Values{}
	form.Add("uname", "testUser")
	form.Add("passwd", "test123")
	form.Add("register", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(RegisterHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	assert.Contains(t, string(body), string(error2.DuplicateUserName))
}

func TestRegisterHandlerGetTemplate(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(RegisterHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, string(body), "Register")
}

func TestValidateInput(t *testing.T) {
	assert.True(t, validateInput("test_1", "test123"))
	assert.False(t, validateInput("", ""))
	assert.False(t, validateInput("test?", "test123"))
}

func TestWrapperValidCookie(t *testing.T) {
	setup()
	defer after()
	// create Session
	sessionToken, expires := CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodGet, "/updateCalendar", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response := httptest.NewRecorder()
	time.Sleep(3 + time.Second)
	Wrapper(LoginHandler).ServeHTTP(response, request)
	cookies := response.Result().Cookies()[0]
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	assert.Equal(t, "testUser", getUsernameBySessionToken(cookies.Value))
	assert.Equal(t, "session_token", cookies.Name)
	assert.Less(t, expires, cookies.Expires)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Calendar")
}

func TestWrapperInvalidCookie(t *testing.T) {
	setup()
	defer after()
	// create Session
	CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodGet, "/updateCalendar", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "wrong_value",
	})
	response := httptest.NewRecorder()
	Wrapper(LoginHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	body, _ := io.ReadAll(response.Result().Body)
	assert.Contains(t, string(body), "Error")
}

func TestGetUsernameBySessionToken(t *testing.T) {
	setup()
	defer after()
	// create Session
	sessionToken, _ := CreateSession("testUser")
	username := getUsernameBySessionToken(sessionToken)
	assert.Equal(t, "testUser", username)

	InitServer()
	username = getUsernameBySessionToken(sessionToken)
	assert.Equal(t, "", username)
}

func TestGetUserBySessionTokenSuccessful(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	sessionToken, _ := CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	user, err := GetUserBySessionToken(request)
	assert.Nil(t, err)
	assert.Equal(t, "testUser", user.UserName)
}

func TestGetUserBySessionTokenUnsuccessfulNoCookie(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	// create Session
	CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	user, err := GetUserBySessionToken(request)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserBySessionTokenUnsuccessfulNoSession(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "value",
	})
	user, err := GetUserBySessionToken(request)
	assert.Error(t, err)
	assert.Equal(t, "cannot get User", err.Error())
	assert.Nil(t, user)
}

func TestGetUserBySessionTokenUnsuccessfulNoUser(t *testing.T) {
	setup()
	defer after()
	// create Session
	sessionToken, _ := CreateSession("testUser")
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	user, err := GetUserBySessionToken(request)
	assert.Error(t, err)
	assert.Equal(t, "cannot get User", err.Error())
	assert.Nil(t, user)
}

func setup() {
	configuration.ReadFlags()
	InitServer()
	templates.Init()
	dataModel.InitDataModel(dataPath)
}
