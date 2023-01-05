package export

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_cal/configuration"
	"go_cal/dataModel"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// an Vorlesung orientiert

func createServer(auth AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		Wrapper(auth,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello client")
			}))
}

func TestWithoutPW(t *testing.T) {
	ts := createServer(func(name, pwd string) bool {
		return true
	})
	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "wrong status")
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t,
		http.StatusText(http.StatusUnauthorized)+"\n",
		string(body), "wrong message")
}

func doRequestWithPassword(t *testing.T, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("<username>", "<password>")
	res, err := client.Do(req)
	assert.NoError(t, err)
	return res
}

func TestWithWrongPW(t *testing.T) {
	var receivedName, receivedPwd string
	ts := createServer(func(name, pwd string) bool {
		receivedName = name
		receivedPwd = pwd
		return false // <--- deny any request
	})
	defer ts.Close()
	res := doRequestWithPassword(t, ts.URL)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "wrong status")
	assert.Equal(t, "<username>", receivedName, "wrong username")
	assert.Equal(t, "<password>", receivedPwd, "wrong password")
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t,
		http.StatusText(http.StatusUnauthorized)+"\n",
		string(body), "wrong message")
}

func TestWithCorrectPW(t *testing.T) {
	var receivedName, receivedPwd string
	ts := createServer(func(name, pwd string) bool {
		receivedName = name
		receivedPwd = pwd
		return true // <--- accept any request
	})
	defer ts.Close()
	res := doRequestWithPassword(t, ts.URL)
	assert.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
	assert.Equal(t, "<username>", receivedName, "wrong username")
	assert.Equal(t, "<password>", receivedPwd, "wrong password")
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello client\n", string(body), "wrong message")
}

func TestCheckUserValid(t *testing.T) {
	setup()
	defer after()
	_, err := dataModel.Dm.AddUser("testUser", "test", 1)
	assert.Nil(t, err)
	assert.True(t, CheckUserValid("testUser", "test"))
	assert.False(t, CheckUserValid("testUser", "wrongPassword"))
}

func setup() {
	configuration.ReadFlags()
	dataModel.InitDataModel(dataPath)
}
