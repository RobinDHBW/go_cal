package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	uTest := NewUser("test1", "test", 1, 3)
	assert.EqualValues(t, 1, uTest.Id)
	assert.EqualValues(t, "test1", uTest.UserName)
	assert.EqualValues(t, "test", uTest.Password)
	assert.EqualValues(t, 3, uTest.UserLevel)
}

func TestNewAppointment(t *testing.T) {
	uTest := NewUser("test1", "test", 1, 3)
	aTest := NewAppointment("test", "hello 123", time.Date(2022, 12, 25, 11, 11, 11, 11, time.UTC), uTest.Id, false, 0, false, "")
	assert.EqualValues(t, "test", aTest.Title)
	assert.EqualValues(t, "hello 123", aTest.Description)
	assert.EqualValues(t, "2022-12-25T11:11:11Z", aTest.DateTime.Format(time.RFC3339))
	assert.EqualValues(t, 1, aTest.Userid)
	assert.EqualValues(t, false, aTest.Timeseries.Repeat)
	assert.EqualValues(t, 0, aTest.Timeseries.Intervall)
	assert.EqualValues(t, false, aTest.Share.Public)
	assert.EqualValues(t, "", aTest.Share.Url)

}
