package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/authentication"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestTerminEditHandler_InvalidRequest(t *testing.T) {
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
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
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
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}

func TestTerminEditHandler_editTermin(t *testing.T) {
	defer after()
	authentication.InitServer()
	dataModel.InitDataModel("../data/test")
	user, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession()
	form := url.Values{}
	form.Add("editTermin", "x")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	addAppointments(user.Id)

	request = initValidSession()
	form = url.Values{}
	form.Add("editTermin", "1")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	request = initValidSession()
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
	authentication.InitServer()
	dataModel.InitDataModel("../data/test")
	user, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession()
	form := url.Values{}
	form.Add("editTerminSubmit", "x")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	addAppointments(user.Id)

	request = initValidSession()
	form = url.Values{}
	form.Add("editTerminSubmit", "1")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	request = initValidSession()
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
	authentication.InitServer()
	dataModel.InitDataModel("../data/test")
	user, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession()
	form := url.Values{}
	form.Add("deleteTerminSubmit", "x")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminEditHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	addAppointments(user.Id)

	request = initValidSession()
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
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	addAppointments(user.Id)
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
	assert.Equal(t, 3, appIndex, "index 2 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 3)
	assert.Equal(t, 4, appIndex, "index 3 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 0)
	assert.Equal(t, 5, appIndex, "index 0 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 1)
	assert.Equal(t, 2, appIndex, "index 1 test")
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
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form := url.Values{}
	form.Add("dateBegin", "keinDatum")
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	request.PostForm = form
	err := EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.InvalidInput, request.Host+"/listTermin"), err)

	// TODO: http und localhost
	request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form = url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("keinDatum"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	request.PostForm = form
	err = EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.InvalidInput, request.Host+"/listTermin"), err)

	// TODO: http und localhost
	request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form = url.Values{}
	form.Add("dateBegin", tThen.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tNow.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	request.PostForm = form
	err = EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.CreateError(error2.EndBeforeBegin, request.Host+"/listTermin"), err)

	// TODO: http und localhost
	request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
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
	dataModel.InitDataModel("../data/test")
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	//apIDstart := dataModel.GetApID()

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form := url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	form.Add("chooseRepeat", "week")
	request.PostForm = form

	err := EditTerminFromInput(request, false, user, 1)
	assert.Equal(t, error2.DisplayedError{}, err)
	assert.Equal(t, "TestTitel", user.Appointments[1].Title)
	assert.Equal(t, "TestContent", user.Appointments[1].Description)
	assert.Equal(t, tNow.Format("2006-01-02T15:04"), user.Appointments[1].DateTimeStart.Format("2006-01-02T15:04"))
	assert.Equal(t, tThen.Format("2006-01-02T15:04"), user.Appointments[1].DateTimeEnd.Format("2006-01-02T15:04"))
	assert.Equal(t, 7, user.Appointments[1].Timeseries.Intervall)
	assert.Equal(t, true, user.Appointments[1].Timeseries.Repeat)
	assert.Equal(t, user.Id, user.Appointments[1].Userid)

}

func TestEditTerminFromInputCorrectInputEdit(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)

	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
	form := url.Values{}
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	form.Add("chooseRepeat", "week")
	request.PostForm = form

	dataModel.Dm.AddAppointment(user.Id, "t", "c", "here", time.Now(), time.Now(), false, 0, true)

	err := EditTerminFromInput(request, true, user, 1)
	assert.Equal(t, error2.DisplayedError{}, err)
	assert.Equal(t, "TestTitel", user.Appointments[1].Title)
	assert.Equal(t, "TestContent", user.Appointments[1].Description)
	assert.Equal(t, tNow.Format("2006-01-02T15:04"), user.Appointments[1].DateTimeStart.Format("2006-01-02T15:04"))
	assert.Equal(t, tThen.Format("2006-01-02T15:04"), user.Appointments[1].DateTimeEnd.Format("2006-01-02T15:04"))
	assert.Equal(t, 7, user.Appointments[1].Timeseries.Intervall)
	assert.Equal(t, true, user.Appointments[1].Timeseries.Repeat)
	assert.Equal(t, user.Id, user.Appointments[1].Userid)

	addAppointments(user.Id)
}

func after() {
	_ = os.RemoveAll("../data/test/")
	_ = os.MkdirAll("../data/test/", 777)
}

func initValidSession() *http.Request {
	templates.Init()

	sessionToken, _ := authentication.CreateSession("testUser")
	cookieValue, _ := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
	// TODO: http und localhost
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/editTermin", nil)

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
