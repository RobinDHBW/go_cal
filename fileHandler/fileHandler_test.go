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
	fH.SyncToFile(write, user.Id)

	fString := fH.ReadFromFile(1)

	var rUser User
	json.Unmarshal([]byte(fString), &rUser)
	return rUser
}

func TestFileHandler_SyncToFile(t *testing.T) {
	user := NewUser("test", "test", 1)
	fH := NewFH("../data/test")
	rUser := fileWriteRead(user, fH)

	assert.EqualValues(t, user, rUser)
}

func TestFileHandler_ReadFromFile(t *testing.T) {

	user := NewUser("test", "test", 1)
	fH := NewFH("../data/test")
	rUser := fileWriteRead(user, fH)

	assert.EqualValues(t, user, rUser)

}
