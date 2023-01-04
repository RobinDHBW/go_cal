package calendar

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestUpdateCalendarHandler_InvalidRequest(t *testing.T) {
	defer after()
	authentication.InitServer()
	templates.Init()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("Testuser", "test", 0)
	assert.Nil(t, err)

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/updateCalendar", nil)
	form := url.Values{}
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	expString, _ := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
	assert.Equal(t, expString, response.Result().Cookies()[0].Value)

	request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/updateCalendar", nil)
	cookieValue, _ := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
	request.AddCookie(&http.Cookie{
		Name:  "fe_parameter",
		Value: cookieValue,
	})
	form = url.Values{}
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
}

func TestUpdateCalendarHandler_CalendarButtons(t *testing.T) {
	defer after()
	authentication.InitServer()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 0)
	assert.Nil(t, err)

	request := initValidSession("updateCalendar")
	form := url.Values{}
	form.Add("next", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	fv := frontendHandling.FrontendView{
		Month:         time.Now().Month(),
		Year:          time.Now().Year(),
		TerminPerSite: 7,
		TerminSite:    1,
		MinDate:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 00, 0, 0, time.UTC),
	}
	fv.NextMonth()
	expString, _ := frontendHandling.GetFeCookieString(fv)
	assert.Equal(t, expString, response.Result().Cookies()[0].Value)

	request = initValidSession("updateCalendar")
	form = url.Values{}
	form.Add("prev", "")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	fv = frontendHandling.FrontendView{
		Month:         time.Now().Month(),
		Year:          time.Now().Year(),
		TerminPerSite: 7,
		TerminSite:    1,
		MinDate:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 00, 0, 0, time.UTC),
	}
	fv.PrevMonth()
	expString, _ = frontendHandling.GetFeCookieString(fv)
	assert.Equal(t, expString, response.Result().Cookies()[0].Value)

	request = initValidSession("updateCalendar")
	fv = frontendHandling.FrontendView{
		Month:         time.Now().Month(),
		Year:          time.Now().Year(),
		TerminPerSite: 7,
		TerminSite:    1,
		MinDate:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 00, 0, 0, time.UTC),
	}

	form = url.Values{}
	form.Add("today", "")
	cookieValue, _ := frontendHandling.GetFeCookieString(fv)
	request.AddCookie(&http.Cookie{
		Name:  "fe_parameter",
		Value: cookieValue,
	})
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	fv.CurrentMonth()
	expString, _ = frontendHandling.GetFeCookieString(fv)
	assert.Equal(t, expString, response.Result().Cookies()[0].Value)
}

func TestUpdateCalendarHandler_ChooseMonth(t *testing.T) {
	defer after()
	authentication.InitServer()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 0)
	assert.Nil(t, err)

	request := initValidSession("updateCalendar")
	form := url.Values{}
	form.Add("choose", "")
	form.Add("chooseYear", "2022")
	form.Add("chooseMonth", "12")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	fv := frontendHandling.FrontendView{
		Month:         12,
		Year:          2022,
		TerminPerSite: 7,
		TerminSite:    1,
		MinDate:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 00, 0, 0, time.UTC),
	}
	expString, _ := frontendHandling.GetFeCookieString(fv)
	assert.Equal(t, expString, response.Result().Cookies()[0].Value)

	request = initValidSession("updateCalendar")
	form = url.Values{}
	form.Add("choose", "")
	form.Add("chooseYear", "kein Jahr")
	form.Add("chooseMonth", "12")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	request = initValidSession("updateCalendar")
	form = url.Values{}
	form.Add("choose", "")
	form.Add("chooseYear", "2022")
	form.Add("chooseMonth", "kein Monat")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	request = initValidSession("updateCalendar")
	form = url.Values{}
	form.Add("choose", "")
	form.Add("chooseYear", "2022")
	form.Add("chooseMonth", "13")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(UpdateCalendarHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}

func after() {
	os.RemoveAll("../data/test/")
	os.MkdirAll("../data/test/", 777)
}

func initValidSession(path string) *http.Request {
	templates.Init()

	sessionToken, _ := authentication.CreateSession("testUser")
	cookieValue, _ := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/"+path, nil)

	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	request.AddCookie(&http.Cookie{
		Name:  "fe_parameter",
		Value: cookieValue,
	})
	return request
}
