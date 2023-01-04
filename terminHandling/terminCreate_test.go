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
	"time"
)

func TestTerminCreateHandler_InvalidRequest(t *testing.T) {
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
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
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
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
}

func TestTerminCreateHandler_CreateTermin(t *testing.T) {
	templates.Init()
	authentication.InitServer()
	defer after()
	dataModel.InitDataModel("../data/test")
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)

	request := initValidSession("createTermin")
	form := url.Values{}
	form.Add("createTermin", "")
	request.PostForm = form
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	request = initValidSession("createTermin")
	form = url.Values{}
	form.Add("createTerminSubmit", "")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	request = initValidSession("createTermin")
	form = url.Values{}
	form.Add("createTerminSubmit", "")
	form.Add("dateBegin", tNow.Format("2006-01-02T15:04"))
	form.Add("dateEnd", tThen.Format("2006-01-02T15:04"))
	form.Add("repeat", "week")
	form.Add("title", "TestTitel")
	form.Add("content", "TestContent")
	form.Add("chooseRepeat", "week")
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusFound, response.Result().StatusCode)
	locationHeader, err := response.Result().Location()
	assert.NoError(t, err)
	assert.Equal(t, "/listTermin", locationHeader.Path)

	request = initValidSession("createTermin")
	form = url.Values{}
	request.PostForm = form
	response = httptest.NewRecorder()
	http.HandlerFunc(TerminCreateHandler).ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
