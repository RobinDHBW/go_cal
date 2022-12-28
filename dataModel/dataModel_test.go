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

var uList = []data.User{data.NewUser("test1", "test", 1, 3), data.NewUser("test2", "test", 2, 2), data.NewUser("test3", "test", 3, 0)}

func fileWriteRead(user data.User, fH *fileHandler.FileHandler) data.User {
	write, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fH.SyncToFile(write, user.Id)

	fString := fH.ReadFromFile(1)

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
	os.MkdirAll("../data/test/", 777)
}

func TestNewDM(t *testing.T) {
	dataPath := "../data/test"

	fH := fileHandler.NewFH("../data/test")
	for _, uD := range uList {
		fileWriteRead(uD, &fH)
	}
	dataModel := NewDM(dataPath)

	defer after()

	//Check if dataPath correct and UserList correct
	assert.EqualValues(t, uList, dataModel.UserList)
}

func TestDataModel_GetUserById(t *testing.T) {
	dataPath := "../data/test"

	fH := fileHandler.NewFH("../data/test")
	for _, uD := range uList {
		fileWriteRead(uD, &fH)
	}
	dataModel := NewDM(dataPath)

	defer after()

	uID := 1
	user := dataModel.GetUserById(uID)

	assert.EqualValues(t, uID, user.Id)
}

func TestDataModel_AddUser(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	user := dataModel.AddUser("test", "abc", 1, nil)
	userFile := dataModel.fH.ReadFromFile(user.Id)
	var user2 data.User

	json.Unmarshal([]byte(userFile), &user2)

	//test if user has same attributes
	//test if file on disk has same attributes

	assert.EqualValues(t, "test", user.UserName)
	assert.EqualValues(t, 1, user.UserLevel)
	assert.EqualValues(t, 0, len(user.Appointments))
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user.Password))

	assert.EqualValues(t, "test", user2.UserName)
	assert.EqualValues(t, 1, user2.UserLevel)
	assert.EqualValues(t, 0, len(user2.Appointments))
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user2.Password))

}

func TestDataModel_AddAppointment(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	tNow := time.Now()
	user := dataModel.AddUser("test", "abc", 1, nil)
	user = dataModel.AddAppointment(user.Id, data.NewAppointment("test", "hello123", tNow, user.Id, false, 0, false, ""))

	assert.EqualValues(t, "test", user.Appointments[0].Title)
	assert.EqualValues(t, tNow, user.Appointments[0].DateTime)
	assert.EqualValues(t, user.Id, user.Appointments[0].Userid)
	assert.EqualValues(t, false, user.Appointments[0].Share.Public)
	assert.EqualValues(t, false, user.Appointments[0].Timeseries.Repeat)

}

func TestDataModel_ComparePW(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	defer after()

	user := dataModel.AddUser("test", "abc", 1, nil)
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user.Password))
	assert.EqualValues(t, false, dataModel.ComparePW("123", user.Password))
}
