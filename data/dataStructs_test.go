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
	aTest := NewAppointment("test", "hello 123", "here", time.Date(2022, 12, 25, 11, 11, 11, 11, time.UTC), time.Date(2022, 12, 25, 12, 12, 12, 12, time.UTC), 0, uTest.Id, false, 0, false)
	assert.EqualValues(t, "test", aTest.Title)
	assert.EqualValues(t, "hello 123", aTest.Description)
	assert.EqualValues(t, "2022-12-25T11:11:11Z", aTest.DateTimeStart.Format(time.RFC3339))
	assert.EqualValues(t, "2022-12-25T12:12:12Z", aTest.DateTimeEnd.Format(time.RFC3339))
	assert.EqualValues(t, 1, aTest.Userid)
	assert.EqualValues(t, false, aTest.Timeseries.Repeat)
	assert.EqualValues(t, 0, aTest.Timeseries.Intervall)
	assert.EqualValues(t, false, aTest.Share.Public)
}

func TestAppointment_GetDescriptionFromInterval(t *testing.T) {
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	ap := NewAppointment("test", "hallo 123", "here", tNow, tThen, 1, 1, true, 1, false)
	assert.Equal(t, "täglich", ap.GetDescriptionFromInterval())
	ap.Timeseries.Intervall = 7
	assert.Equal(t, "wöchentlich", ap.GetDescriptionFromInterval())
	ap.Timeseries.Intervall = 365
	assert.Equal(t, "jährlich", ap.GetDescriptionFromInterval())
	ap.Timeseries.Intervall = 30
	assert.Equal(t, "monatlich", ap.GetDescriptionFromInterval())
	ap.Timeseries.Intervall = -5
	assert.Equal(t, "keine", ap.GetDescriptionFromInterval())

}

func TestShare_GetUsernameFromUrl(t *testing.T) {
	text := "http://localhost:8080/terminVoting?invitor=test&termin=test&token=jWAgIWSYiPxDiauBNPfQusername=Testuser"
	tNow := time.Now()
	tThen := tNow.Add(time.Hour * time.Duration(1))
	ap := NewAppointment("test", "hallo 123", "here", tNow, tThen, 1, 1, true, 1, false)
	result := ap.Share.GetUsernameFromUrl(text)
	assert.Equal(t, "", result)
	text = "http://localhost:8080/terminVoting?invitor=test&termin=test&token=jWAgIWSYiPxDiauBNPfQ&username=Testuser"
	result = ap.Share.GetUsernameFromUrl(text)
	assert.Equal(t, "Testuser", result)
}
