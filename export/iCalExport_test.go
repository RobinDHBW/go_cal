package export

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_cal/dataModel"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewVEvent(t *testing.T) {
	stamp := time.Now()
	subject := NewVEvent("test1", "Here", "test", "test", "Test", stamp, stamp, stamp)

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
	dataPath := "../data/test"
	dataModel := dataModel.NewDM(dataPath)

	defer os.RemoveAll(dataPath)
	user, err := dataModel.AddUser("test", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	t1 := time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC)
	t1End := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)

	user, ap := dataModel.AddAppointment(user.Id, "test", "search for", "here", t1, t1End, false, 0, false)
	subject := NewICal(ap)

	assert.EqualValues(t, fmt.Sprintf("%d", ap.Id), subject.VEvent.UID)
	assert.EqualValues(t, "here", subject.VEvent.Location)
	assert.EqualValues(t, "test", subject.VEvent.Summary)
	assert.EqualValues(t, "search for", subject.VEvent.Description)
	assert.EqualValues(t, "PUBLIC", subject.VEvent.Class)
	assert.EqualValues(t, "1", subject.VEvent.UID)
	assert.EqualValues(t, t1, subject.VEvent.DTStart)
	assert.EqualValues(t, t1End, subject.VEvent.DTEnd)

}

func TestICal_ToString(t *testing.T) {
	dataPath := "../data/test"
	dataModel := dataModel.NewDM(dataPath)

	defer os.RemoveAll(dataPath)
	user, err := dataModel.AddUser("test", "abc", 1)
	if err != nil {
		t.FailNow()
	}

	t1 := time.Date(2022, 12, 24, 10, 00, 00, 00, time.UTC)
	t1End := time.Date(2022, 12, 24, 11, 00, 00, 00, time.UTC)

	user, ap := dataModel.AddAppointment(user.Id, "test", "search for", "here", t1, t1End, false, 0, false)
	subject := NewICal(ap)
	check := subject.ToString()

	splits := strings.Split(check, "\n")
	assert.EqualValues(t, "BEGIN:VCALENDAR", splits[0])
	assert.EqualValues(t, "END:VCALENDAR", splits[len(splits)-1])
}
