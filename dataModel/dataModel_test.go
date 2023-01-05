package dataModel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_cal/configuration"
	"go_cal/data"
	"go_cal/fileHandler"
	"go_cal/templates"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// var uList = []data.User{data.NewUser("test1", "test", 1, 3), data.NewUser("test2", "test", 2, 2), data.NewUser("test3", "test", 3, 0)}
var uMap = map[int]data.User{1: data.NewUser("test1", "test", 1, 3), 2: data.NewUser("test2", "test", 2, 2), 3: data.NewUser("test3", "test", 3, 0)}

const dataPath = "../data/test/DM"

func after() {
	err := os.RemoveAll(dataPath)
	if err != nil {
		return
	}
}

func fileWriteRead(user data.User, fH *fileHandler.FileHandler, t *testing.T) data.User {
	write, err := json.Marshal(user)
	if err != nil {
		t.FailNow()
	}
	//defer after()
	fH.SyncToFile(write, user.Id)

	fString := fH.ReadFromFile(user.Id)

	var rUser data.User
	err = json.Unmarshal([]byte(fString), &rUser)
	if err != nil {
		t.FailNow()
	}
	return rUser
}

func TestNewDM(t *testing.T) {
	//dataPath := "../data/test"

	fH := fileHandler.NewFH(dataPath)
	for _, uD := range uMap {
		fileWriteRead(uD, &fH, t)
	}

	dataModel := NewDM(dataPath)
	defer after()

	//Check if dataPath correct and UserList correct
	assert.EqualValues(t, uMap[0], dataModel.UserMap[0])
}

func TestDataModel_GetUserById(t *testing.T) {
	//dataPath := "../data/test"

	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test1", "abc", 1)
	if err != nil {
		t.FailNow()
	}
	uID := user.Id
	user = dataModel.GetUserById(uID)

	assert.EqualValues(t, uID, user.Id)
}

func TestDataModel_GetUserByName(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test2", "abc", 1)

	if err != nil {
		t.FailNow()
	}

	user = dataModel.GetUserByName("test2")
	assert.EqualValues(t, "test2", user.UserName)

	user = dataModel.GetUserByName("test3")
	assert.Nil(t, user)
	user = Dm.GetUserByName("te")
	assert.Nil(t, user)
}

func TestDataModel_AddUser(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	fH := fileHandler.NewFH(dataPath)
	defer after()

	user, err := dataModel.AddUser("test3", "abc", 1)
	assert.Nil(t, err)

	userFile := fH.ReadFromFile(user.Id)
	var user2 data.User

	err = json.Unmarshal([]byte(userFile), &user2)
	if err != nil {
		t.FailNow()
	}

	//test if user has same attributes
	//test if file on disk has same attributes

	assert.EqualValues(t, "test3", user.UserName)
	assert.EqualValues(t, 1, user.UserLevel)
	assert.EqualValues(t, 0, len(user.Appointments))
	assert.EqualValues(t, true, Dm.ComparePW("abc", user.Password))

	assert.EqualValues(t, "test3", user2.UserName)
	assert.EqualValues(t, 1, user2.UserLevel)
	assert.EqualValues(t, 0, len(user2.Appointments))
	assert.EqualValues(t, true, Dm.ComparePW("abc", user2.Password))

	user, err = dataModel.AddUser("test3", "abc", 1)
	assert.Error(t, err)

}

func TestDataModel_AddAppointment(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	defer after()

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	user, err := dataModel.AddUser("test4", "abc", 1)
	if err != nil {
		t.FailNow()
	}
	user, ap := dataModel.AddAppointment(user.Id, "test", "hello123", "here", tNow, tThen, false, 0, false)

	assert.EqualValues(t, "test", ap.Title)
	assert.EqualValues(t, tNow, ap.DateTimeStart)
	assert.EqualValues(t, tThen, ap.DateTimeEnd)
	assert.EqualValues(t, user.Id, ap.Userid)
	assert.EqualValues(t, false, ap.Share.Public)
	assert.EqualValues(t, false, ap.Timeseries.Repeat)
}

func TestDataModel_DeleteAppointment(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test5", "abc", 1)

	if err != nil {
		t.FailNow()
	}

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	user, ap := dataModel.AddAppointment(user.Id, "test", "hello 123", "here", tNow, tThen, false, 0, false)
	user, ap = dataModel.AddAppointment(user.Id, "test1", "hello 123", "here", tNow, tThen, false, 0, false)
	user, ap = dataModel.AddAppointment(user.Id, "test2", "hello 123", "here", tNow, tThen, false, 0, false)

	lenAp := len(user.Appointments)
	user = dataModel.DeleteAppointment(ap.Id, user.Id)

	_, ok := user.Appointments[ap.Id]

	assert.EqualValues(t, lenAp-1, len(user.Appointments))
	assert.False(t, ok)
}

func TestDataModel_EditAppointment(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test6", "abc", 1)

	if err != nil {
		t.FailNow()
	}

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	title := "test"
	//ap1 := data.NewAppointment()

	user, ap := dataModel.AddAppointment(user.Id, title, "hello 123", "here", tNow, tThen, false, 0, false)

	assert.EqualValues(t, title, ap.Title)

	title = "test123"
	ap.Title = title
	user = dataModel.EditAppointment(user.Id, ap)

	assert.EqualValues(t, title, user.Appointments[ap.Id].Title)
}

func TestDataModel_GetAppointmentsBySearchString(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test7", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	t1 := time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC)
	t1End := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)

	user, _ = dataModel.AddAppointment(user.Id, "test", "search for", "here", t1, t1End, false, 0, false)
	user, _ = dataModel.AddAppointment(user.Id, "test1", "catch me if you can", "here", t1, t1End, false, 0, false)
	user, _ = dataModel.AddAppointment(user.Id, "test2", "qwertzuiopasdfghjklyxcvbnm123456789", "Here", t1, t1End, false, 0, false)

	//user = dataModel.AddAppointment(dataModel.AddAppointment(dataModel.AddAppointment(user.Id, ap1).Id, ap2).Id, ap3)

	_, check := dataModel.GetAppointmentsBySearchString(user.Id, "test")
	assert.EqualValues(t, len(user.Appointments), len(*check))

	_, check = dataModel.GetAppointmentsBySearchString(user.Id, "catch")
	assert.EqualValues(t, 1, len(*check))

	_, check = dataModel.GetAppointmentsBySearchString(user.Id, "123456")
	assert.EqualValues(t, 1, len(*check))

}

func TestDataModel_ComparePW(t *testing.T) {
	//dataPath := "../data/test"
	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test8", "abc", 1)
	if err != nil {
		t.FailNow()
	}
	assert.EqualValues(t, true, Dm.ComparePW("abc", user.Password))
	assert.EqualValues(t, false, Dm.ComparePW("123", user.Password))
}

func TestDataModel_SetVotingForTokenSuccessful(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// zweiter Terminvorschlag
	beginDate, _ = time.Parse("2006-01-02T15:04", "2023-02-10T22:00")
	endDate, _ = time.Parse("2006-01-02T15:04", "2023-02-10T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	assert.Equal(t, 2, len(user.SharedAppointments["Terminfindung1"]))

	// user einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Anna", "Terminfindung1", "Peter"), "Anna")
	assert.Nil(t, err)
	// zweiten User einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Hans", "Terminfindung1", "Peter"), "Hans")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(user.SharedAppointments["Terminfindung1"][0].Share.Tokens))
	assert.Equal(t, 2, len(user.SharedAppointments["Terminfindung1"][0].Share.Voting))
	// initial haben beide User bei beiden Terminen abgesagt
	assert.False(t, user.SharedAppointments["Terminfindung1"][0].Share.Voting[0])
	assert.False(t, user.SharedAppointments["Terminfindung1"][0].Share.Voting[1])
	assert.False(t, user.SharedAppointments["Terminfindung1"][1].Share.Voting[0])
	assert.False(t, user.SharedAppointments["Terminfindung1"][1].Share.Voting[1])
	// URLs für die 2 Termine müssen pro User gleich sein
	assert.Equal(t, user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0], user.SharedAppointments["Terminfindung1"][1].Share.Tokens[0])
	assert.Equal(t, user.SharedAppointments["Terminfindung1"][0].Share.Tokens[1], user.SharedAppointments["Terminfindung1"][1].Share.Tokens[1])
	// token extrahieren
	urlAnna, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	urlHans, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[1])
	assert.Nil(t, err)
	tokenAnna := urlAnna.Query().Get("token")
	tokenHans := urlHans.Query().Get("token")
	// Hans sagt für den ersten Termin zu, für den zweiten ab
	keys := make([]int, 0)
	keys = append(keys, 0)
	err = Dm.SetVotingForToken(user, keys, "Terminfindung1", tokenHans, "Hans")
	assert.Nil(t, err)
	// Hans hat für den ersten Termin zugesagt
	assert.True(t, user.SharedAppointments["Terminfindung1"][0].Share.Voting[1])
	// und für den 2. abgesagt
	assert.False(t, user.SharedAppointments["Terminfindung1"][1].Share.Voting[1])
	// Anna sagt für beide Termine zu
	keys = make([]int, 0)
	keys = append(keys, 0, 1)
	err = Dm.SetVotingForToken(user, keys, "Terminfindung1", tokenAnna, "Anna")
	assert.Nil(t, err)
	// Anna hat für beide Termin zugesagt
	assert.True(t, user.SharedAppointments["Terminfindung1"][0].Share.Voting[0])
	assert.True(t, user.SharedAppointments["Terminfindung1"][1].Share.Voting[0])
}

func TestDataModel_SetVotingForTokenUnsuccessful(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)

	// user einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Anna", "Terminfindung1", "Peter"), "Anna")
	assert.Nil(t, err)

	// token extrahieren
	urlAnna, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	tokenAnna := urlAnna.Query().Get("token")
	// Anna sagt für den einzigen Termin zu
	keys := make([]int, 0)
	keys = append(keys, 0)
	// falscher Titel
	err = Dm.SetVotingForToken(user, keys, "falsche Terminfindung", tokenAnna, "Anna")
	assert.Error(t, err)
}

func TestDataModel_IsVotingAllowedSuccessful(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Anna", "Terminfindung1", "Peter"), "Anna")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	assert.True(t, IsVotingAllowed("Terminfindung1", token, user, "Anna"))
}

func TestDataModel_IsVotingAllowedUnsuccessfulNoUser(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Anna", "Terminfindung1", "Peter"), "Anna")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	assert.False(t, IsVotingAllowed("Terminfindung1", token, nil, "Anna"))
}

func TestDataModel_IsVotingAllowedUnsuccessfulNoApp(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Anna", "Terminfindung1", "Peter"), "Anna")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	assert.False(t, IsVotingAllowed("nicht vorhandene Terminfindung", token, user, "Anna"))
}

func TestDataModel_IsVotingAllowedUnsuccessfulWrongQuery(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	// user einladen
	err = Dm.AddTokenToSharedAppointment(user.Id, "Terminfindung1", CreateURL("Anna", "Terminfindung1", "Peter"), "Anna")
	assert.Nil(t, err)
	// token extrahieren
	tokenUrl, err := url.Parse(user.SharedAppointments["Terminfindung1"][0].Share.Tokens[0])
	assert.Nil(t, err)
	token := tokenUrl.Query().Get("token")
	assert.False(t, IsVotingAllowed("Terminfindung1", token, user, "Hans"))
}

func TestCreateURL(t *testing.T) {
	createdUrl := CreateURL("Peter", "Terminvorschlag1", "Hans")
	assert.Contains(t, createdUrl, "/terminVoting?invitor=Hans&termin=Terminvorschlag1&token=")
	assert.Contains(t, createdUrl, "&username=Peter")
}

func TestCreateToken(t *testing.T) {
	InitSeed()
	token1 := createToken(20)
	assert.NotEqual(t, "", token1)
	token2 := createToken(20)
	assert.NotEqual(t, "", token2)
	assert.NotEqual(t, token1, token2)
}

func TestDataModel_GetAppointmentsForUser(t *testing.T) {
	setup()
	defer after()

	user, err := Dm.AddUser("Peter", "test123", 1)
	assert.Nil(t, err)
	// Terminfindung erstellen
	beginDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T22:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2023-01-03T23:00")
	Dm.AddSharedAppointment(user.Id, "Terminfindung1", "here", beginDate, endDate, false, 0, true)
	assert.Equal(t, &user.Appointments, Dm.GetAppointmentsForUser(user.Id))
}

func setup() {
	configuration.ReadFlags()
	dir, _ := os.Getwd()
	templates.Init(filepath.Join(dir, ".."))
	InitDataModel(dataPath)
}

func TestDataModel_AddSharedAppointment(t *testing.T) {
	dataModel := NewDM(dataPath)
	defer after()

	user, err := dataModel.AddUser("test5", "abc", 1)
	if err != nil {
		t.FailNow()
	}
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	apIDbefore := apID
	lenBefore := len(user.SharedAppointments["test"])
	dataModel.AddSharedAppointment(user.Id, "test", "here", tNow, tThen, false, 0, false)

	assert.Equal(t, apIDbefore+1, apID)
	assert.Equal(t, lenBefore+1, len(user.SharedAppointments["test"]))
	expAp := data.NewAppointment("test", "", "here", tNow, tThen, apIDbefore, user.Id, false, 0, false)
	assert.EqualValues(t, expAp.Title, user.SharedAppointments["test"][lenBefore].Title)
	assert.EqualValues(t, expAp.Location, user.SharedAppointments["test"][lenBefore].Location)
	assert.EqualValues(t, expAp.DateTimeEnd, user.SharedAppointments["test"][lenBefore].DateTimeEnd)
	assert.EqualValues(t, expAp.DateTimeStart, user.SharedAppointments["test"][lenBefore].DateTimeStart)
	assert.EqualValues(t, expAp.Timeseries, user.SharedAppointments["test"][lenBefore].Timeseries)
	assert.EqualValues(t, expAp.Share, user.SharedAppointments["test"][lenBefore].Share)

	_ = dataModel.AddTokenToSharedAppointment(user.Id, "test", "TestURL", "invitedUser")
	_ = dataModel.AddTokenToSharedAppointment(user.Id, "test", "TestURL2", "invitedUser2")
	dataModel.AddSharedAppointment(user.Id, "test", "here", tNow.Add(time.Hour*time.Duration(1)), tThen.Add(time.Hour*time.Duration(1)), false, 0, false)
	assert.Equal(t, &user.SharedAppointments["test"][lenBefore+1].Share.Tokens, &user.SharedAppointments["test"][lenBefore].Share.Tokens)
	assert.Equal(t, true, user.SharedAppointments["test"][lenBefore+1].Share.Public)
	assert.Equal(t, false, user.SharedAppointments["test"][lenBefore+1].Share.Voting[0])
	assert.Equal(t, false, user.SharedAppointments["test"][lenBefore+1].Share.Voting[1])

	dataModel.AddSharedAppointment(user.Id, "test1", "here", tNow, tThen, false, 0, false)
	assert.Equal(t, 2, len(user.SharedAppointments["test"]))
	assert.Equal(t, 1, len(user.SharedAppointments["test1"]))
}

func TestDataModel_AddTokenToSharedAppointment(t *testing.T) {
	InitDataModel(dataPath)
	defer after()

	user, err := Dm.AddUser("test6", "abc", 1)
	if err != nil {
		log.Fatal("error not nil")
		t.FailNow()
	}
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	lenBefore := len(user.SharedAppointments["test"])
	Dm.AddSharedAppointment(user.Id, "test", "here", tNow, tThen, false, 0, false)

	err = Dm.AddTokenToSharedAppointment(user.Id, "test", "/xyz?username=invitedUser", "invitedUser")
	assert.Nil(t, err)
	err = Dm.AddTokenToSharedAppointment(user.Id, "test", "/xyz?username=invitedUser", "invitedUser")
	assert.NotNil(t, err)
	err = Dm.AddTokenToSharedAppointment(user.Id, "test", "/xyz?username=invitedUser2", "invitedUser2")
	assert.Nil(t, err)

	assert.Equal(t, "/xyz?username=invitedUser", user.SharedAppointments["test"][lenBefore].Share.Tokens[0])
	assert.Equal(t, "/xyz?username=invitedUser2", user.SharedAppointments["test"][lenBefore].Share.Tokens[1])

	assert.Equal(t, lenBefore+2, len(user.SharedAppointments["test"][lenBefore].Share.Tokens))
	assert.Equal(t, false, user.SharedAppointments["test"][lenBefore].Share.Voting[0])
	assert.Equal(t, false, user.SharedAppointments["test"][lenBefore].Share.Voting[1])

}

func TestDataModel_DeleteSharedAppointment(t *testing.T) {
	InitDataModel(dataPath)
	defer after()

	user, err := Dm.AddUser("test7", "abc", 1)

	if err != nil {
		t.FailNow()
	}

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))

	Dm.AddSharedAppointment(user.Id, "test", "here", tNow, tThen, false, 0, false)
	Dm.AddSharedAppointment(user.Id, "test", "here", tNow.Add(time.Hour*time.Duration(1)), tThen.Add(time.Hour*time.Duration(1)), false, 0, false)
	Dm.AddSharedAppointment(user.Id, "test1", "here", tNow, tThen, false, 0, false)

	assert.Equal(t, 2, len(user.SharedAppointments["test"]))
	assert.Equal(t, 1, len(user.SharedAppointments["test1"]))

	user = Dm.DeleteSharedAppointment("test", user.Id)

	_, ok := user.SharedAppointments["test"]

	assert.Equal(t, 0, len(user.SharedAppointments["test"]))
	assert.Equal(t, 1, len(user.SharedAppointments["test1"]))
	assert.False(t, ok)
}
