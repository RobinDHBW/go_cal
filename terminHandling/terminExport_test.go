package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/dataModel"
	"go_cal/export"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const dataPath = "../data/test/ICAL"

func createServer(auth export.AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		export.Wrapper(auth,
			func(w http.ResponseWriter, r *http.Request) {
				ICalHandler(w, r)
			}))
}

func TestICalHandler(t *testing.T) {
	dataModel.InitDataModel(dataPath)
	user, err := dataModel.Dm.AddUser("test", "abc", 3)
	if err != nil {
		t.FailNow()
	}

	defer after()

	//var _, _ string
	ts := createServer(func(name, pwd string) bool {
		//receivedName = name
		//receivedPwd = pwd
		return true // <--- accept any request
	})
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("test", "abc")
	res, err := client.Do(req)
	assert.NoError(t, err)
	body, err := io.ReadAll(res.Body)

	//create mock
	ics := export.NewICal(dataModel.Dm.GetAppointmentsForUser(user.Id)).ToString()
	assert.EqualValues(t, []byte(ics), body)

}
