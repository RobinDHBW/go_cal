package dataModel

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"go_cal/fileHandler"
	"testing"
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

func init() {
	fH := fileHandler.NewFH("../data/test")
	for _, uD := range uList {
		fileWriteRead(uD, &fH)
	}
}

func TestNewDM(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	//Check if dataPath correct and UserList correct
	assert.EqualValues(t, uList, dataModel.UserList)
}

func TestDataModel_GetUserById(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	uID := 1
	user := dataModel.GetUserById(uID)

	assert.EqualValues(t, uID, user.Id)
}

func TestDataModel_AddUser(t *testing.T) {
	dataPath := "../data/test"
	dataModel := NewDM(dataPath)

	user := dataModel.AddUser("test", "abc", 1, nil)
	//test if user has same attributes
	//test if file on disk has same attributes

	assert.EqualValues(t, "test", user.UserName)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, 0, len(user.Appointments))
	assert.EqualValues(t, true, dataModel.ComparePW("abc", user.Password))
}
