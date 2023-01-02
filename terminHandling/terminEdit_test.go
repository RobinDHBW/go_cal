package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/dataModel"
	"go_cal/frontendHandling"
	"os"
	"testing"
	"time"
)

func TestGetTerminFromEditIndex(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	addAppointments(user.Id)
	// Termine ab 01.01.2023
	fv := frontendHandling.FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2023, 1, 1, 11, 11, 1, 1, time.Local),
	}
	// Exspected order of appointmentIds: 4 1 2 3
	appIndex := GetTerminFromEditIndex(*user, fv, 2)
	assert.Equal(t, 2, appIndex, "index 2 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 3)
	assert.Equal(t, 3, appIndex, "index 3 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 0)
	assert.Equal(t, 4, appIndex, "index 0 test")

	appIndex = GetTerminFromEditIndex(*user, fv, 1)
	assert.Equal(t, 1, appIndex, "index 1 test")
}

func TestGetRepeatingMode(t *testing.T) {
	mode := "none"
	assert.Equal(t, 0, GetRepeatingMode(mode))

	mode = "day"
	assert.Equal(t, 1, GetRepeatingMode(mode))

	mode = "week"
	assert.Equal(t, 7, GetRepeatingMode(mode))

	mode = "month"
	assert.Equal(t, 30, GetRepeatingMode(mode))

	mode = "year"
	assert.Equal(t, 365, GetRepeatingMode(mode))

	mode = "other"
	assert.Equal(t, 0, GetRepeatingMode(mode))
}

func TestEditTerminFromInputIncorrectInput(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	addAppointments(user.Id)
}

func TestEditTerminFromInputCorrectInputCreate(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	addAppointments(user.Id)
}

func TestEditTerminFromInputCorrectInputEdit(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	addAppointments(user.Id)
}

func after() {
	os.RemoveAll("../data/test/")
	os.MkdirAll("../data/test/", 777)
}
