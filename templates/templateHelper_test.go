// Matrikelnummern:
// 9495107, 4706893, 9608900

package templates

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	assert.Nil(t, TempInit)
	assert.Nil(t, TempLogin)
	assert.Nil(t, TempRegister)
	assert.Nil(t, TempError)
	assert.Nil(t, TempTerminList)
	assert.Nil(t, TempTerminEdit)
	assert.Nil(t, TempCreateTermin)
	assert.Nil(t, TempShareTermin)
	assert.Nil(t, TempCreateShareTermin)
	assert.Nil(t, TempEditShareTermin)
	assert.Nil(t, TempTerminVoting)
	assert.Nil(t, TempTerminVotingSuccess)
	assert.Nil(t, TempSearchTermin)
	dir, _ := os.Getwd()
	Init(filepath.Join(dir, ".."))
	assert.NotNil(t, TempInit)
	assert.NotNil(t, TempLogin)
	assert.NotNil(t, TempRegister)
	assert.NotNil(t, TempError)
	assert.NotNil(t, TempTerminList)
	assert.NotNil(t, TempTerminEdit)
	assert.NotNil(t, TempCreateTermin)
	assert.NotNil(t, TempShareTermin)
	assert.NotNil(t, TempCreateShareTermin)
	assert.NotNil(t, TempEditShareTermin)
	assert.NotNil(t, TempTerminVoting)
	assert.NotNil(t, TempTerminVotingSuccess)
	assert.NotNil(t, TempSearchTermin)
}
