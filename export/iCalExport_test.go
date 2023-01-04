package export

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"go_cal/dataModel"
	"os"
	"strings"
	"testing"
	"time"
)

const dataPath = "../data/test/IC"

func after() {
	err := os.RemoveAll(dataPath)
	if err != nil {
		return
	}
}

func TestNewVEvent(t *testing.T) {
	stamp := time.Now()
	subject := NewVEvent("test1", "Here", "test", "test", "Test", stamp, stamp, stamp, data.TimeSeries{false, 0})

	assert.EqualValues(t, "test1", subject.UID)
	assert.EqualValues(t, "Here", subject.Location)
	assert.EqualValues(t, "test", subject.Summary)
	assert.EqualValues(t, "test", subject.Description)
	assert.EqualValues(t, "Test", subject.Class)
	assert.EqualValues(t, "test1", subject.UID)
	assert.EqualValues(t, stamp, subject.DTStart)
	assert.EqualValues(t, stamp, subject.DTEnd)
	assert.EqualValues(t, stamp, subject.DTStamp)
}

func TestNewICal(t *testing.T) {

	dM := dataModel.NewDM(dataPath)

	defer after()
	user, err := dM.AddUser("ical", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	t1 := time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC)
	t1End := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)

	user, ap := dM.AddAppointment(user.Id, "test", "search for", "here", t1, t1End, false, 0, false)
	subject := NewICal(dM.GetAppointmentsForUser(user.Id))

	assert.EqualValues(t, fmt.Sprintf("%d", ap.Id), subject.VEvent[0].UID)
	assert.EqualValues(t, "here", subject.VEvent[0].Location)
	assert.EqualValues(t, "test", subject.VEvent[0].Summary)
	assert.EqualValues(t, "search for", subject.VEvent[0].Description)
	assert.EqualValues(t, "PUBLIC", subject.VEvent[0].Class)
	assert.EqualValues(t, "1", subject.VEvent[0].UID)
	assert.EqualValues(t, t1, subject.VEvent[0].DTStart)
	assert.EqualValues(t, t1End, subject.VEvent[0].DTEnd)

}

func TestICal_ToString(t *testing.T) {
	//dataPath := "../data/test"
	dM := dataModel.NewDM(dataPath)

	defer after()
	user, err := dM.AddUser("test", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	t1 := time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC)
	t1End := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)

	dM.AddAppointment(user.Id, "test", "search for", "here", t1, t1End, true, 1, false)
	subject := NewICal(dM.GetAppointmentsForUser(user.Id))
	check := subject.ToString()

	splits := strings.Split(check, "\n")
	assert.EqualValues(t, "BEGIN:VCALENDAR", splits[0])
	assert.EqualValues(t, "END:VCALENDAR", splits[len(splits)-1])

	index := 0
	for i, val := range splits {
		if strings.Contains(val, "RRULE") {
			index = i
		}

	}

	assert.True(t, strings.Contains(splits[index], "DAILY"))
}
