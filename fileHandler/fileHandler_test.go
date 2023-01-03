package fileHandler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"log"
	"os"
	"strconv"

	//"go_cal/dataModel"
	"testing"
)

const dataPath = "../data/test/FH"

func fileWriteRead(user data.User, fH *FileHandler) data.User {
	write, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	fH.SyncToFile(write, user.Id)

	fString := fH.ReadFromFile(user.Id)

	var rUser data.User
	err = json.Unmarshal([]byte(fString), &rUser)
	if err != nil {
		log.Fatal(err)
	}
	return rUser
}

func after() {
	err := os.RemoveAll(dataPath)
	if err != nil {
		return
	}
}

func TestNewFH(t *testing.T) {
	//dP := "../data/test"
	fH := NewFH(dataPath)
	defer after()

	assert.EqualValues(t, dataPath, fH.dataPath)
}

func TestFileHandler_SyncToFile(t *testing.T) {
	user := data.NewUser("test", "test", 1, 3)
	fH := NewFH(dataPath)
	rUser := fileWriteRead(user, &fH)

	defer after()

	assert.EqualValues(t, user, rUser)
}

func TestFileHandler_ReadFromFile(t *testing.T) {

	user := data.NewUser("test", "test", 1, 3)
	fH := NewFH(dataPath)
	rUser := fileWriteRead(user, &fH)

	defer after()

	assert.EqualValues(t, user, rUser)
}

func TestFileHandler_ReadAll(t *testing.T) {
	uList := make([]data.User, 0)
	for i := 0; i < 10; i++ {
		uList = append(uList, data.NewUser("test"+strconv.Itoa(i), "test", i, 3))
	}

	fH := NewFH(dataPath)
	for _, uD := range uList {
		fileWriteRead(uD, &fH)
	}
	defer after()

	var rUList []data.User
	for _, uString := range fH.ReadAll() {
		var user data.User
		err := json.Unmarshal([]byte(uString), &user)
		if err != nil {
			t.FailNow()
		}

		rUList = append(rUList, user)
	}
	assert.EqualValues(t, uList, rUList)
}
