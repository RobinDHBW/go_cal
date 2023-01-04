package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/data"
	"go_cal/dataModel"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

//var App1 data.Appointment
//var App2 data.Appointment
//var App3 data.Appointment
//var App4 data.Appointment
//var App5 data.Appointment

func TestTerminHandler_InvalidRequest(t *testing.T) {
	templates.Init()
	authentication.InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request, _ := http.NewRequest(http.MethodPost, "/editTermin", nil)
	form := url.Values{}
	request.PostForm = form
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "cookie123",
	})
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)

	sessionToken, _ := authentication.CreateSession("testUser")

	request, _ = http.NewRequest(http.MethodPost, "/editTermin", nil)
	form = url.Values{}
	request.PostForm = form
	request.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}

func TestTerminHandler_ButtonBackAndSearch(t *testing.T) {
	templates.Init()
	authentication.InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession("listTermin")
	form := url.Values{}
	form.Add("calendarBack", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	request = initValidSession("listTermin")
	form = url.Values{}
	form.Add("terminlistBack", "")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	request = initValidSession("listTermin")
	form = url.Values{}
	form.Add("searchTerminSubmit", "")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	request = initValidSession("listTermin")
	form = url.Values{}
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestTerminHandler_SubmitTermin(t *testing.T) {
	templates.Init()
	authentication.InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession("listTermin")
	form := url.Values{}
	form.Add("submitTermin", "")
	form.Add("numberPerSite", "keineZahl")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	request = initValidSession("listTermin")
	form = url.Values{}
	form.Add("submitTermin", "")
	form.Add("numberPerSite", "5")
	form.Add("siteChoose", "keineZahl")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	request = initValidSession("listTermin")
	form = url.Values{}
	form.Add("submitTermin", "")
	form.Add("numberPerSite", "5")
	form.Add("siteChoose", "1")
	form.Add("dateChoose", "keinDatum")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	request = initValidSession("listTermin")
	form = url.Values{}
	form.Add("submitTermin", "")
	form.Add("numberPerSite", "5")
	form.Add("siteChoose", "1")
	form.Add("dateChoose", "2022-12-12")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminHandler).ServeHTTP(response, request)
	cookie := response.Result().Cookies()[0]
	fvExp := frontendHandling.FrontendView{
		Month:         time.Now().Month(),
		Year:          time.Now().Year(),
		TerminPerSite: 5,
		TerminSite:    1,
		MinDate:       time.Date(2022, 12, 12, 0, 0, 0, 0, time.UTC),
	}
	expCookie, _ := frontendHandling.GetFeCookieString(fvExp)
	assert.Equal(t, expCookie, cookie.Value)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

}

func TestGetFirstTerminOfRepeatingInDate(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	ap1, ap2, ap3, ap4, ap5 := addAppointments(user.Id)
	// Termine ab 01.01.2023
	fv := frontendHandling.FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2023, 1, 1, 11, 11, 1, 1, time.Local),
	}

	// Check single Appointment
	resApp := frontendHandling.GetFirstTerminOfRepeatingInDate(*ap1, fv)
	assert.Equal(t, resApp, *ap1, "single test")

	// Check daily
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(*ap5, fv)
	ap5.DateTimeStart = time.Date(2023, 01, 01, 15, 0, 0, 0, time.Local)
	ap5.DateTimeEnd = time.Date(2023, 01, 01, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, *ap5, "daily test")

	// Check weekly
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(*ap2, fv)
	ap2.DateTimeStart = time.Date(2023, 01, 02, 14, 0, 0, 0, time.Local)
	ap2.DateTimeEnd = time.Date(2023, 01, 02, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, *ap2, "weekly test")

	// Check monthly
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(*ap3, fv)
	ap3.DateTimeStart = time.Date(2023, 01, 11, 15, 0, 0, 0, time.Local)
	ap3.DateTimeEnd = time.Date(2023, 01, 11, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, *ap3, "monthly test")

	// Check yearly
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(*ap4, fv)
	ap4.DateTimeStart = time.Date(2023, 12, 12, 15, 0, 0, 0, time.Local)
	ap4.DateTimeEnd = time.Date(2023, 12, 12, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, *ap4, "yearly test")
}

func addAppointments(id int) (ap1, ap2, ap3, ap4, ap5 *data.Appointment) {
	// Einzelner Termin: 7 Dez 2022
	_, ap1 = dataModel.Dm.AddAppointment(id, "titel1", "beschreibung1", "here", time.Date(2022, 12, 7, 12, 0, 0, 0, time.Local),
		time.Date(2022, 12, 7, 14, 0, 0, 0, time.Local), false, 0, false)
	// Wöchentlicher Termin: ab 12 Dez 2022
	_, ap2 = dataModel.Dm.AddAppointment(id, "titel2", "beschreibung2", "here", time.Date(2022, 12, 12, 14, 0, 0, 0, time.Local),
		time.Date(2022, 12, 12, 17, 0, 0, 0, time.Local), true, 7, false)
	// Monatlicher Termin: ab 11 Nov 2022
	_, ap3 = dataModel.Dm.AddAppointment(id, "titel3", "beschreibung3", "here", time.Date(2022, 11, 11, 15, 0, 0, 0, time.Local),
		time.Date(2022, 11, 11, 17, 0, 0, 0, time.Local), true, 30, false)
	// Jährlicher Termin ab 12 Dez 2021
	_, ap4 = dataModel.Dm.AddAppointment(id, "titel4", "beschreibung4", "here", time.Date(2021, 12, 12, 15, 0, 0, 0, time.Local),
		time.Date(2021, 12, 12, 17, 0, 0, 0, time.Local), true, 365, false)
	// Täglicher Termin ab 17 Dez 2022
	_, ap5 = dataModel.Dm.AddAppointment(id, "titel5", "beschreibung5", "here", time.Date(2022, 12, 17, 15, 0, 0, 0, time.Local),
		time.Date(2022, 12, 17, 17, 0, 0, 0, time.Local), true, 1, false)
	return
}
