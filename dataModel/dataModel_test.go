package dataModel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"go_cal/fileHandler"
	"os"
	"testing"
	"time"
)

// var uList = []data.User{data.NewUser("test1", "test", 1, 3), data.NewUser("test2", "test", 2, 2), data.NewUser("test3", "test", 3, 0)}
var uMap = map[int]data.User{1: data.NewUser("test1", "test", 1, 3), 2: data.NewUser("test2", "test", 2, 2), 3: data.NewUser("test3", "test", 3, 0)}

func fileWriteRead(user data.User, fH *fileHandler.FileHandler) data.User {
	write, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fH.SyncToFile(write, user.Id)

	fString := fH.ReadFromFile(user.Id)

	var rUser data.User
	json.Unmarshal([]byte(fString), &rUser)
	return rUser
}

//func init() {
//	//fH := fileHandler.NewFH("../data/test")
//	//for _, uD := range uList {
//	//	fileWriteRead(uD, &fH)
//	//}
//}

func after() {
	os.RemoveAll("../data/test/")
	//os.MkdirAll("../data/test/", 777)
}

func TestNewDM(t *testing.T) {
	dataPath := "../data/test"

	fH := fileHandler.NewFH("../data/test")
	for _, uD := range uMap {
		fileWriteRead(uD, &fH)
	}
	dataModel := NewDM(dataPath)

	defer after()

	//Check if dataPath correct and UserList correct
	assert.EqualValues(t, uMap, dataModel.UserMap)
}

func TestDataModel_GetUserById(t *testing.T) {
	dataPath := "../data/test"

	fH := fileHandler.NewFH("../data/test")
	for _, uD := range uMap {
		fileWriteRead(uD, &fH)
	}
	dataModel := NewDM(dataPath)

	defer after()

	uID := 1
	user := dataModel.GetUserById(uID)

	assert.EqualValues(t, uID, user.Id)
}

func TestDataModel_GetUserByName(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	user, err := dataModel.AddUser("test", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	user = dataModel.GetUserByName("test")
	assert.EqualValues(t, "test", user.UserName)

	user = dataModel.GetUserByName("test1")
	assert.Nil(t, user)
	user = dataModel.GetUserByName("te")
	assert.Nil(t, user)
}

func TestDataModel_AddUser(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	user, err := dataModel.AddUser("test", "abc", 1)

	assert.Nil(t, err)

	userFile := dataModel.fH.ReadFromFile(user.Id)
	var user2 data.User

	json.Unmarshal([]byte(userFile), &user2)

	//test if user has same attributes
	//test if file on disk has same attributes

	assert.EqualValues(t, "test", user.UserName)
	assert.EqualValues(t, "test", user.UserName)
	assert.EqualValues(t, 1, user.UserLevel)
	assert.EqualValues(t, 0, len(user.Appointments))
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user.Password))

	assert.EqualValues(t, "test", user2.UserName)
	assert.EqualValues(t, 1, user2.UserLevel)
	assert.EqualValues(t, 0, len(user2.Appointments))
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user2.Password))

	user, err = dataModel.AddUser("test", "abc", 1)
	assert.Error(t, err)

}

func TestDataModel_AddAppointment(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	user, err := dataModel.AddUser("test", "abc", 1)
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
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()
	user, err := dataModel.AddUser("test", "abc", 1)
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
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()
	user, err := dataModel.AddUser("test", "abc", 1)
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

func TestDataModel_GetAppointmentByTimeFrame(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()
	user, err := dataModel.AddUser("test", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	t1 := time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC)
	t1End := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)

	t2 := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)
	t2End := time.Date(2022, 12, 24, 12, 00, 00, 00, time.UTC)

	t3 := time.Date(2022, 12, 24, 12, 00, 00, 00, time.UTC)
	t3End := time.Date(2022, 12, 24, 13, 00, 00, 00, time.UTC)

	user, _ = dataModel.AddAppointment(user.Id, "test", "hello 123", "here", t1, t1End, false, 0, false)
	user, _ = dataModel.AddAppointment(user.Id, "test1", "hello 123", "here", t2, t2End, false, 0, false)
	user, _ = dataModel.AddAppointment(user.Id, "test2", "hello 123", "Here", t3, t3End, false, 0, false)

	//user = dataModel.AddAppointment(dataModel.AddAppointment(dataModel.AddAppointment(user.Id, ap1).Id, ap2).Id, ap3)

	_, check := dataModel.GetAppointmentsByTimeFrame(user.Id, time.Date(2022, 12, 24, 9, 59, 00, 00, time.UTC), time.Date(2022, 12, 24, 13, 01, 00, 00, time.UTC))
	assert.EqualValues(t, len(user.Appointments), len(*check))

	_, check2 := dataModel.GetAppointmentsByTimeFrame(user.Id, time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC), time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC))
	assert.EqualValues(t, len(user.Appointments)-1, len(*check2))

}

func TestDataModel_GetAppointmentsBySearchString(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()
	user, err := dataModel.AddUser("test", "abc", 1)
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
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	user, err := dataModel.AddUser("test", "abc", 1)
	if err != nil {
		t.FailNow()
	}
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user.Password))
	assert.EqualValues(t, false, dataModel.ComparePW("123", user.Password))
}
