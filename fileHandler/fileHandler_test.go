package fileHandler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"os"

	//"go_cal/dataModel"
	"testing"
)

func fileWriteRead(user data.User, fH *FileHandler) data.User {
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

func after() {
	os.RemoveAll("../data/test/")
	os.MkdirAll("../data/test/", 777)
}

func TestNewFH(t *testing.T) {
	dP := "../data/test"
	fH := NewFH(dP)

	assert.EqualValues(t, dP, fH.dataPath)
}

func TestFileHandler_SyncToFile(t *testing.T) {
	user := data.NewUser("test", "test", 1, 3)
	fH := NewFH("../data/test")
	rUser := fileWriteRead(user, &fH)

	defer after()

	assert.EqualValues(t, user, rUser)
}

func TestFileHandler_ReadFromFile(t *testing.T) {

	user := data.NewUser("test", "test", 1, 3)
	fH := NewFH("../data/test")
	rUser := fileWriteRead(user, &fH)

	defer after()

	assert.EqualValues(t, user, rUser)
}

func TestFileHandler_ReadAll(t *testing.T) {
	uList := []data.User{data.NewUser("test1", "test", 1, 3), data.NewUser("test2", "test", 2, 3), data.NewUser("test3", "test", 3, 3)}
	fH := NewFH("../data/test")
	for _, uD := range uList {
		fileWriteRead(uD, &fH)
	}
	//defer after()

	//rSList := fH.ReadAll()
	rUList := []data.User{}
	for _, uString := range fH.ReadAll() {
		var user data.User
		json.Unmarshal([]byte(uString), &user)

		rUList = append(rUList, user)
	}
	assert.EqualValues(t, uList, rUList)
}
