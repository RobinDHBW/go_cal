package terminHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"go_cal/dataModel"
	"go_cal/frontendHandling"
	"os"
	"strconv"
	"testing"
	"time"
)

var App1 data.Appointment
var App2 data.Appointment
var App3 data.Appointment
var App4 data.Appointment
var App5 data.Appointment

func TestGetFirstTerminOfRepeatingInDate(t *testing.T) {
	dataModel.InitDataModel()
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

	// Check single Appointment
	resApp := frontendHandling.GetFirstTerminOfRepeatingInDate(App1, fv)
	assert.Equal(t, resApp, App1, "single test")

	// Check daily
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(App5, fv)
	App5.DateTimeStart = time.Date(2023, 01, 01, 15, 0, 0, 0, time.Local)
	App5.DateTimeEnd = time.Date(2023, 01, 01, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, App5, "daily test")

	// Check weekly
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(App2, fv)
	App2.DateTimeStart = time.Date(2023, 01, 02, 14, 0, 0, 0, time.Local)
	App2.DateTimeEnd = time.Date(2023, 01, 02, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, App2, "weekly test")

	// Check monthly
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(App3, fv)
	App3.DateTimeStart = time.Date(2023, 01, 11, 15, 0, 0, 0, time.Local)
	App3.DateTimeEnd = time.Date(2023, 01, 11, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, App3, "monthly test")

	// Check yearly
	resApp = frontendHandling.GetFirstTerminOfRepeatingInDate(App4, fv)
	App4.DateTimeStart = time.Date(2023, 12, 12, 15, 0, 0, 0, time.Local)
	App4.DateTimeEnd = time.Date(2023, 12, 12, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, App4, "yearly test")

	_ = os.Remove("../files/" + strconv.FormatInt(int64(user.Id), 10) + ".json")

}

func addAppointments(id int) {
	// Einzelner Termin: 7 Dez 2022
	App1 = data.NewAppointment("titel1", "beschreibung1",
		time.Date(2022, 12, 7, 12, 0, 0, 0, time.Local),
		time.Date(2022, 12, 7, 14, 0, 0, 0, time.Local),
		id, false, 0, false, "")
	// Wöchentlicher Termin: ab 12 Dez 2022
	App2 = data.NewAppointment("titel2", "beschreibung2",
		time.Date(2022, 12, 12, 14, 0, 0, 0, time.Local),
		time.Date(2022, 12, 12, 17, 0, 0, 0, time.Local),
		id, true, 7, false, "")
	// Monatlicher Termin: ab 11 Nov 2022
	App3 = data.NewAppointment("titel3", "beschreibung3",
		time.Date(2022, 11, 11, 15, 0, 0, 0, time.Local),
		time.Date(2022, 11, 11, 17, 0, 0, 0, time.Local),
		id, true, 30, false, "")
	// Jährlicher Termin ab 12 Dez 2021
	App4 = data.NewAppointment("titel4", "beschreibung4",
		time.Date(2021, 12, 12, 15, 0, 0, 0, time.Local),
		time.Date(2021, 12, 12, 17, 0, 0, 0, time.Local),
		id, true, 365, false, "")
	// Täglicher Termin ab 17 Dez 2022
	App5 = data.NewAppointment("titel5", "beschreibung5",
		time.Date(2022, 12, 17, 15, 0, 0, 0, time.Local),
		time.Date(2022, 12, 17, 17, 0, 0, 0, time.Local),
		id, true, 1, false, "")
	dataModel.Dm.AddAppointment(id, App1)
	dataModel.Dm.AddAppointment(id, App2)
	dataModel.Dm.AddAppointment(id, App3)
	dataModel.Dm.AddAppointment(id, App4)
	dataModel.Dm.AddAppointment(id, App5)

}
