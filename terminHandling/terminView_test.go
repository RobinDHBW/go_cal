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

var app1 data.Appointment
var app2 data.Appointment
var app3 data.Appointment
var app4 data.Appointment
var app5 data.Appointment

func TestGetTerminList(t *testing.T) {
	dataModel.InitDataModel()
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	addAppointments(user.Id)
	// Termine ab 11.12.2022
	fv := frontendHandling.FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2022, 11, 11, 11, 11, 1, 1, time.Local),
	}
	apps := GetTerminList(user.Appointments, fv)
	// Richtige Reihenfolge der Termine
	exp := make([]data.Appointment, 0, 1)
	exp = append(exp, app3)
	exp = append(exp, app1)
	exp = append(exp, app2)
	expApp4 := app4
	expApp4.DateTimeStart = expApp4.DateTimeStart.AddDate(1, 0, 0)
	expApp4.DateTimeEnd = expApp4.DateTimeEnd.AddDate(1, 0, 0)
	exp = append(exp, expApp4)
	exp = append(exp, app5)
	assert.Equal(t, apps, exp, "test 1 equal")

	// Termine ab 20.12.2022
	fv = frontendHandling.FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2022, 12, 20, 11, 11, 1, 1, time.Local),
	}
	apps = GetTerminList(user.Appointments, fv)
	// Richtige Reihenfolge der Termine
	exp = make([]data.Appointment, 0, 1)
	expApp5 := app5
	expApp5.DateTimeStart = expApp5.DateTimeStart.AddDate(0, 0, 3)
	expApp5.DateTimeEnd = expApp5.DateTimeEnd.AddDate(0, 0, 3)
	exp = append(exp, expApp5)
	expApp2 := app2
	expApp2.DateTimeStart = expApp2.DateTimeStart.AddDate(0, 0, 14)
	expApp2.DateTimeEnd = expApp2.DateTimeEnd.AddDate(0, 0, 14)
	exp = append(exp, expApp2)
	expApp3 := app3
	expApp3.DateTimeStart = expApp3.DateTimeStart.AddDate(0, 2, 0)
	expApp3.DateTimeEnd = expApp3.DateTimeEnd.AddDate(0, 2, 0)
	exp = append(exp, expApp3)
	expApp4 = app4
	expApp4.DateTimeStart = expApp4.DateTimeStart.AddDate(2, 0, 0)
	expApp4.DateTimeEnd = expApp4.DateTimeEnd.AddDate(2, 0, 0)
	exp = append(exp, expApp4)
	assert.Equal(t, apps, exp, "test 2 equal")

	_ = os.Remove("../files/" + strconv.FormatInt(int64(user.Id), 10) + ".json")
}

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

	// Check single appointment
	resApp := GetFirstTerminOfRepeatingInDate(app1, fv)
	assert.Equal(t, resApp, app1, "single test")

	// Check daily
	resApp = GetFirstTerminOfRepeatingInDate(app5, fv)
	app5.DateTimeStart = time.Date(2023, 01, 01, 15, 0, 0, 0, time.Local)
	app5.DateTimeEnd = time.Date(2023, 01, 01, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, app5, "daily test")

	// Check weekly
	resApp = GetFirstTerminOfRepeatingInDate(app2, fv)
	app2.DateTimeStart = time.Date(2023, 01, 02, 14, 0, 0, 0, time.Local)
	app2.DateTimeEnd = time.Date(2023, 01, 02, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, app2, "weekly test")

	// Check monthly
	resApp = GetFirstTerminOfRepeatingInDate(app3, fv)
	app3.DateTimeStart = time.Date(2023, 01, 11, 15, 0, 0, 0, time.Local)
	app3.DateTimeEnd = time.Date(2023, 01, 11, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, app3, "monthly test")

	// Check yearly
	resApp = GetFirstTerminOfRepeatingInDate(app4, fv)
	app4.DateTimeStart = time.Date(2023, 12, 12, 15, 0, 0, 0, time.Local)
	app4.DateTimeEnd = time.Date(2023, 12, 12, 17, 0, 0, 0, time.Local)
	assert.Equal(t, resApp, app4, "yearly test")

	_ = os.Remove("../files/" + strconv.FormatInt(int64(user.Id), 10) + ".json")

}

func addAppointments(id int) {
	// Einzelner Termin: 7 Dez 2022
	app1 = data.NewAppointment("titel1", "beschreibung1",
		time.Date(2022, 12, 7, 12, 0, 0, 0, time.Local),
		time.Date(2022, 12, 7, 14, 0, 0, 0, time.Local),
		id, false, 0, false, "")
	// Wöchentlicher Termin: ab 12 Dez 2022
	app2 = data.NewAppointment("titel2", "beschreibung2",
		time.Date(2022, 12, 12, 14, 0, 0, 0, time.Local),
		time.Date(2022, 12, 12, 17, 0, 0, 0, time.Local),
		id, true, 7, false, "")
	// Monatlicher Termin: ab 11 Nov 2022
	app3 = data.NewAppointment("titel3", "beschreibung3",
		time.Date(2022, 11, 11, 15, 0, 0, 0, time.Local),
		time.Date(2022, 11, 11, 17, 0, 0, 0, time.Local),
		id, true, 30, false, "")
	// Jährlicher Termin ab 12 Dez 2021
	app4 = data.NewAppointment("titel4", "beschreibung4",
		time.Date(2021, 12, 12, 15, 0, 0, 0, time.Local),
		time.Date(2021, 12, 12, 17, 0, 0, 0, time.Local),
		id, true, 365, false, "")
	// Täglicher Termin ab 17 Dez 2022
	app5 = data.NewAppointment("titel5", "beschreibung5",
		time.Date(2022, 12, 17, 15, 0, 0, 0, time.Local),
		time.Date(2022, 12, 17, 17, 0, 0, 0, time.Local),
		id, true, 1, false, "")
	dataModel.Dm.AddAppointment(id, app1)
	dataModel.Dm.AddAppointment(id, app2)
	dataModel.Dm.AddAppointment(id, app3)
	dataModel.Dm.AddAppointment(id, app4)
	dataModel.Dm.AddAppointment(id, app5)

}
