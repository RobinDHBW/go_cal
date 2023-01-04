package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestTerminEditHandler_InvalidRequest(t *testing.T) {
	defer after()
	setup()
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
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
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
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}

func TestTerminEditHandler_editTermin(t *testing.T) {
	defer after()
	setup()
	user, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession("editTermin")
	form := url.Values{}
	form.Add("editTermin", "x")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	addAppointments(user.Id)

	request = initValidSession("editTermin")
	form = url.Values{}
	form.Add("editTermin", "1")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	request = initValidSession("editTermin")
	form = url.Values{}
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listTermin", locationHeader.Path)
}

func TestTerminEditHandler_editTerminSubmit(t *testing.T) {
	defer after()
	setup()
	user, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession("editTermin")
	form := url.Values{}
	form.Add("editTerminSubmit", "x")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	addAppointments(user.Id)

	request = initValidSession("editTermin")
	form = url.Values{}
	form.Add("editTerminSubmit", "1")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	request = initValidSession("editTermin")
	form = url.Values{}
	request.PostForm = form
	form.Add("editTerminSubmit", "1")
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	form.Add("chooseRepeat", "week")
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listTermin", locationHeader.Path)
}

func TestTerminEditHandler_deleteTerminSubmit(t *testing.T) {
	defer after()
	setup()
	user, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession("editTermin")
	form := url.Values{}
	form.Add("deleteTerminSubmit", "x")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	addAppointments(user.Id)

	request = initValidSession("editTermin")
	form = url.Values{}
	request.PostForm = form
	form.Add("deleteTerminSubmit", "1")
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listTermin", locationHeader.Path)
}

func TestGetTerminFromEditIndex(t *testing.T) {
	defer after()
	setup()
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	_, ap2, ap3, ap4, ap5 := addAppointments(user.Id)
	// Termine ab 01.01.2023
	fv := frontendHandling.FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2023, 1, 1, 11, 11, 1, 1, time.Local),
	}
	// Expected order of appointmentIds: 4 1 2 3
	appIndex := GetTerminFromEditIndex(*user, fv, 2)
	assert.Equal(t, ap3.Id, appIndex, "index 2 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 3)
	assert.Equal(t, ap4.Id, appIndex, "index 3 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 0)
	assert.Equal(t, ap5.Id, appIndex, "index 0 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 1)
	assert.Equal(t, ap2.Id, appIndex, "index 1 test")
}

func TestGetRepeatingMode(t *testing.T) {
	mode := "none"
	assert.Equal(t, 0, GetRepeatingMode(mode))

	mode = "day"
	assert.Equal(t, 1, GetRepeatingMode(mode))

	mode = "week"
	assert.Equal(t, 7, GetRepeatingMode(mode))

	mode = "month"
	assert.Equal(t, 30, GetRepeatingMode(mode))

	mode = "year"
	assert.Equal(t, 365, GetRepeatingMode(mode))

	mode = "other"
	assert.Equal(t, 0, GetRepeatingMode(mode))
}

func TestEditTerminFromInputIncorrectInput(t *testing.T) {
	defer after()
	setup()
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Add("dateBegin", "keinDatum")
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	request.PostForm = form
	err := EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.InvalidInput, request.Host+"/listTermin"), err)

	request, _ = http.NewRequest(http.MethodPost, "/", nil)
	form = url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("keinDatum"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	request.PostForm = form
	err = EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.InvalidInput, request.Host+"/listTermin"), err)

	request, _ = http.NewRequest(http.MethodPost, "/", nil)
	form = url.Values{}
	form.Add("dateBegin", tThen.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tNow.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	request.PostForm = form
	err = EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.EndBeforeBegin, request.Host+"/listTermin"), err)

	request, _ = http.NewRequest(http.MethodPost, "/", nil)
	form = url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "")
	form.Add("content", "TestContent")
	request.PostForm = form
	err = EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.TitleIsEmpty, request.Host+"/listTermin"), err)
}

func TestEditTerminFromInputCorrectInputCreate(t *testing.T) {
	defer after()
	setup()
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	form.Add("chooseRepeat", "week")
	request.PostForm = form
	err := EditTerminFromInput(request, false, user, 1)
	maxId := 0 // maxId is the newest ID = added ID
	for k := range user.Appointments {
		maxId = k
	}

	assert.Equal(t, error2.DisplayedError{}, err)
	assert.Equal(t, "TestTitel", user.Appointments[maxId].Title)
	assert.Equal(t, "TestContent", user.Appointments[maxId].Description)
	assert.Equal(t, tNow.Format("2006-01-02T15:04"), user.Appointments[maxId].DateTimeStart.Format("2006-01-02T15:04"))
	assert.Equal(t, tThen.Format("2006-01-02T15:04"), user.Appointments[maxId].DateTimeEnd.Format("2006-01-02T15:04"))
	assert.Equal(t, 7, user.Appointments[maxId].Timeseries.Intervall)
	assert.Equal(t, true, user.Appointments[maxId].Timeseries.Repeat)
	assert.Equal(t, user.Id, user.Appointments[maxId].Userid)

}

func TestEditTerminFromInputCorrectInputEdit(t *testing.T) {
	defer after()
	setup()
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	form.Add("chooseRepeat", "week")
	request.PostForm = form

	_, a := dataModel.Dm.AddAppointment(user.Id, "t", "c", "here", time.Now(), time.Now(), false, 0, true)
	id := a.Id
	err := EditTerminFromInput(request, true, user, id)
	assert.Equal(t, error2.DisplayedError{}, err)
	assert.Equal(t, "TestTitel", user.Appointments[id].Title)
	assert.Equal(t, "TestContent", user.Appointments[id].Description)
	assert.Equal(t, tNow.Format("2006-01-02T15:04"), user.Appointments[id].DateTimeStart.Format("2006-01-02T15:04"))
	assert.Equal(t, tThen.Format("2006-01-02T15:04"), user.Appointments[id].DateTimeEnd.Format("2006-01-02T15:04"))
	assert.Equal(t, 7, user.Appointments[id].Timeseries.Intervall)
	assert.Equal(t, true, user.Appointments[id].Timeseries.Repeat)
	assert.Equal(t, user.Id, user.Appointments[id].Userid)
}

func after() {
	_ = os.RemoveAll("../data/test/")
	_ = os.MkdirAll("../data/test/", 777)
}

func initValidSession(path string) *http.Request {
	sessionToken, _ := authentication.CreateSession("testUser")
	cookieValue, _ := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
	request, _ := http.NewRequest(http.MethodPost, "/"+path, nil)

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
