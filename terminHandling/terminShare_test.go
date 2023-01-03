package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestTerminShareHandler(t *testing.T) {

}

func TestCreateSharedTerminSuccessful(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, error2.DisplayedError{}, displayedError)
	assert.Equal(t, 1, len(user.SharedAppointments["Terminfindung"]))
	assert.Equal(t, "Terminfindung", user.SharedAppointments["Terminfindung"][0].Title)
	assert.Equal(t, user.Id, user.SharedAppointments["Terminfindung"][0].Userid)
	parsed, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	assert.Equal(t, parsed, user.SharedAppointments["Terminfindung"][0].DateTimeStart)
	parsed, _ = time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	assert.Equal(t, parsed, user.SharedAppointments["Terminfindung"][0].DateTimeEnd)
}

func TestCreateSharedTerminUnsuccessfulWrongFormatBegin(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	// wrong format in dateBegin
	form.Add("dateBegin", "2023-0103T22:00")
	form.Add("dateEnd", "2023-01-03T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, string(error2.InvalidInput), displayedError.Text)
}

func TestCreateSharedTerminUnsuccessfulWrongFormatEnd(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	// wrong format in dateEnd
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-0103T23:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, string(error2.InvalidInput), displayedError.Text)
}

func TestCreateSharedTerminUnsuccessfulEndBeforeBegin(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	// user erstellen
	user, err := dataModel.Dm.AddUser("otto", "test123", 1)
	assert.Nil(t, err)
	// request vorbereiten
	request, _ := http.NewRequest(http.MethodPost, "/shareTermin", nil)
	form := url.Values{}
	// wrong format in dateEnd
	form.Add("dateBegin", "2023-01-03T22:00")
	form.Add("dateEnd", "2023-01-03T21:00")
	form.Add("chooseRepeat", "none")
	request.PostForm = form
	displayedError := createSharedTermin(request, user, "Terminfindung")
	assert.Equal(t, string(error2.EndBeforeBegin), displayedError.Text)
}

func TestCreateURL(t *testing.T) {
	url := CreateURL("Peter", "Terminvorschlag1", "Hans")
	assert.Contains(t, url, "/terminVoting?invitor=Hans&termin=Terminvorschlag1&token=")
	assert.Contains(t, url, "&username=Peter")
}

func TestValidateInput(t *testing.T) {
	successful := validateInput("")
	assert.False(t, successful)
	successful = validateInput("test?")
	assert.False(t, successful)
	successful = validateInput("Test123_")
	assert.True(t, successful)
}

func TestCreateToken(t *testing.T) {
	InitSeed()
	token1 := createToken(20)
	assert.NotEqual(t, "", token1)
	token2 := createToken(20)
	assert.NotEqual(t, "", token2)
	assert.NotEqual(t, token1, token2)
}
