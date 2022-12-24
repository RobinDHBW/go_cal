package fileHandler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func fileWriteRead(user User, fH FileHandler) User {
	write, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fH.SyncToFile(string(write))

	fString := fH.ReadFromFile(1)

	var rUser User
	json.Unmarshal([]byte(fString), &rUser)
	return rUser
}

// Check to write json struct to file and readback --> struct must be same
func TestFileHandler_SyncToFile(t *testing.T) {
	user := NewUser("test", "test", 1)
	fH := NewFH("../data")
	rUser := fileWriteRead(user, fH)

	assert.EqualValues(t, user, rUser)
}

// Check to read json from file
func TestFileHandler_ReadFromFile(t *testing.T) {
	//Step 1: common way to create a user
	//Step 2: write to disk
	//Step 3: reread --> must be the same

	user := NewUser("test", "test", 1)
	fH := NewFH("../data")
	rUser := fileWriteRead(user, fH)

	assert.EqualValues(t, user, rUser)

}
