package frontendHandling

import (
	"github.com/stretchr/testify/assert"
	"go_cal/data"
	"go_cal/dataModel"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

//var App1 data.Appointment
//var App2 data.Appointment
//var App3 data.Appointment
//var App4 data.Appointment
//var App5 data.Appointment

func TestFrontendView_GetDaysOfMonth(t *testing.T) {
	fv := FrontendView{
		Month: 2,
		Year:  2022,
	}
	res := fv.GetDaysOfMonth()
	// Feb 2021 has 28 days
	assert.Equal(t, len(res), 28, "Feb 2021: Length not equal")
	exp := make([]int, 28)
	for i := range exp {
		exp[i] = i + 1
	}
	assert.Equal(t, exp, res, "Feb 2021: Slices are not equal")

	fv = FrontendView{
		Month: 2,
		Year:  2024,
	}
	res = fv.GetDaysOfMonth()
	// Feb 2024 has 29 days
	exp = make([]int, 29)
	assert.Equal(t, len(res), 29, "Feb 2024: Length not equal")
	for i := range exp {
		exp[i] = i + 1
	}
	assert.Equal(t, exp, res, "Feb 2024: Slices are not equal")
}

func TestFrontendView_GetDaysBeforeMonthBegin(t *testing.T) {
	fv := FrontendView{
		Month: 2,
		Year:  2022,
	}
	res := fv.GetDaysBeforeMonthBegin()
	// Feb 2022 starts on Tue, so 1 day before
	assert.Equal(t, len(res), 1, "Feb 2022: Length not equal")

	fv = FrontendView{
		Month: 8,
		Year:  2022,
	}
	res = fv.GetDaysBeforeMonthBegin()
	// Aug 2022 starts on Mon, so no day before
	assert.Equal(t, len(res), 0, "Aug 2022: Length not equal")

	fv = FrontendView{
		Month: 5,
		Year:  2022,
	}
	res = fv.GetDaysBeforeMonthBegin()
	// May 2022 starts on Sun, so 6 days before
	assert.Equal(t, len(res), 6, "May 2022: Length not equal")
}

func TestFrontendView_NextMonth(t *testing.T) {
	fv := FrontendView{
		Month: 12,
		Year:  2022,
	}
	fv.NextMonth()
	assert.Equal(t, fv.Month, time.Month(1), "Dez 2022: Month not correct")
	assert.Equal(t, fv.Year, 2023, "Dez 2022: Year not correct")
	fv = FrontendView{
		Month: 6,
		Year:  1993,
	}
	fv.NextMonth()
	assert.Equal(t, fv.Month, time.Month(7), "Jun 1993: Month not correct")
	assert.Equal(t, fv.Year, 1993, "Jun 1993: Year not correct")
}

func TestFrontendView_PrevMonth(t *testing.T) {
	fv := FrontendView{
		Month: 1,
		Year:  2022,
	}
	fv.PrevMonth()
	assert.Equal(t, fv.Month, time.Month(12), "Jan 2022: Month not correct")
	assert.Equal(t, fv.Year, 2021, "Jan 2022: Year not correct")
	fv = FrontendView{
		Month: 6,
		Year:  1993,
	}
	fv.PrevMonth()
	assert.Equal(t, fv.Month, time.Month(5), "Jun 1993: Month not correct")
	assert.Equal(t, fv.Year, 1993, "Jun 1993: Year not correct")
}

func TestFrontendView_CurrentMonth(t *testing.T) {
	fv := FrontendView{
		Month: 1,
		Year:  2022,
	}
	fv.CurrentMonth()
	assert.Equal(t, fv.Month, time.Now().Month(), "Jan 2022: Month not correct")
	assert.Equal(t, fv.Year, time.Now().Year(), "Jan 2022: Year not correct")
}

func TestFrontendView_ChooseMonth(t *testing.T) {
	fv := FrontendView{
		Month: 1,
		Year:  2022,
	}
	err := fv.ChooseMonth(2017, 5)
	assert.Equal(t, fv.Month, time.Month(5), "Jan 2022: Month not correct")
	assert.Equal(t, fv.Year, 2017, "Jan 2022: Year not correct")
	assert.Nil(t, err, "Error not nil")

	fv = FrontendView{
		Month: 1,
		Year:  2022,
	}
	err = fv.ChooseMonth(2017, 13)
	assert.Equal(t, fv.Month, time.Month(1), "Jan 2022: Month not correct")
	assert.Equal(t, fv.Year, 2022, "Jan 2022: Year not correct")
	assert.NotNil(t, err, "Error nil")

	fv = FrontendView{
		Month: 1,
		Year:  2022,
	}
	err = fv.ChooseMonth(2017, 0)
	assert.Equal(t, fv.Month, time.Month(1), "Jan 2022: Month not correct")
	assert.Equal(t, fv.Year, 2022, "Jan 2022: Year not correct")
	assert.NotNil(t, err, "Error nil")

	fv = FrontendView{
		Month: 1,
		Year:  2022,
	}
	err = fv.ChooseMonth(-1, 6)
	assert.Equal(t, fv.Month, time.Month(1), "Jan 2022: Month not correct")
	assert.Equal(t, fv.Year, 2022, "Jan 2022: Year not correct")
	assert.NotNil(t, err, "Error nil")
}

func TestFrontendView_GetCurrentDate(t *testing.T) {
	fv := FrontendView{}
	assert.WithinDuration(t, fv.GetCurrentDate(), time.Now(), 1*time.Second, "Times should be equal")
}

func TestFrontendView_GetAppointmentsForMonth(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	fv := FrontendView{
		Year:  2022,
		Month: 12,
	}
	user, _ := dataModel.Dm.AddUser("Testuser", "pw", 0)
	// Einzelner Termin: 7 Dez 2022
	dataModel.Dm.AddAppointment(user.Id, "titel1", "beschreibung1", "here", time.Date(2022, 12, 7, 12, 0, 0, 0, time.Local),
		time.Date(2022, 12, 7, 14, 0, 0, 0, time.Local), false, 0, false)
	// Wöchentlicher Termin: ab 12 Dez 2022
	dataModel.Dm.AddAppointment(user.Id, "titel2", "beschreibung2", "here", time.Date(2022, 12, 12, 15, 0, 0, 0, time.Local),
		time.Date(2022, 12, 12, 17, 0, 0, 0, time.Local), true, 7, false)
	exp := make([]int, 32)
	exp[7] = 1
	exp[12] = 1
	exp[19] = 1
	exp[26] = 1
	apps := fv.GetAppointmentsForMonth(*user)
	assert.Equal(t, apps, exp, "test single appointment and weekly repeat: Slices not equal")

	// Monatlicher Termin: ab 11 Nov 2022
	dataModel.Dm.AddAppointment(user.Id, "titel3", "beschreibung3", "here", time.Date(2022, 11, 11, 15, 0, 0, 0, time.Local),
		time.Date(2022, 11, 11, 17, 0, 0, 0, time.Local), true, 30, false)
	exp[11] = 1
	apps = fv.GetAppointmentsForMonth(*user)
	assert.Equal(t, apps, exp, "test monthly repeat: Slices not equal")

	// Jährlicher Termin: ab 12 Dez 2021
	dataModel.Dm.AddAppointment(user.Id, "titel4", "beschreibung4", "here", time.Date(2021, 12, 12, 15, 0, 0, 0, time.Local),
		time.Date(2021, 12, 12, 17, 0, 0, 0, time.Local), true, 365, false)
	exp[12]++
	apps = fv.GetAppointmentsForMonth(*user)
	assert.Equal(t, apps, exp, "test yearly repeat: Slices not equal")

	// Täglicher Termin: ab 17 Dez 2022
	dataModel.Dm.AddAppointment(user.Id, "titel5", "beschreibung5", "here", time.Date(2022, 12, 17, 15, 0, 0, 0, time.Local),
		time.Date(2022, 12, 17, 17, 0, 0, 0, time.Local), true, 1, false)
	for i := range exp[17:] {
		exp[17+i]++
	}
	apps = fv.GetAppointmentsForMonth(*user)
	assert.Equal(t, apps, exp, "test yearly repeat: Slices not equal")

	// Test previous month
	fv = FrontendView{
		Year:  2022,
		Month: 11,
	}
	exp = make([]int, 32)
	exp[11] = 1
	apps = fv.GetAppointmentsForMonth(*user)
	assert.Equal(t, apps, exp, "test prev month: Slices not equal")

	// Test following month
	fv = FrontendView{
		Year:  2023,
		Month: 1,
	}
	exp = make([]int, 32)
	for i := range exp {
		exp[i]++
	}
	exp[0] = 0
	exp[11]++
	exp[2]++
	exp[9]++
	exp[16]++
	exp[23]++
	exp[30]++
	apps = fv.GetAppointmentsForMonth(*user)
	assert.Equal(t, apps, exp, "test prev month: Slices not equal")
}

func TestGetFrontendParameters(t *testing.T) {
	// valid Cookie
	cookieValue := "{'Month':10,'Year':2022,'TerminPerSite':10,'TerminSite':2,'MinDate':'2022-09-12T10:00:00Z'}"
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, &http.Cookie{Name: "fe_parameter", Value: cookieValue})
	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	fv, err := GetFrontendParameters(request)
	fvExp := &FrontendView{
		Month:         10,
		Year:          2022,
		TerminPerSite: 10,
		TerminSite:    2,
		MinDate:       time.Date(2022, 9, 12, 10, 00, 0, 0, time.UTC),
	}
	assert.Equal(t, fv, fvExp, "structs should be equal")
	assert.Nil(t, err)

	// no Cookie
	request = &http.Request{}
	fv, err = GetFrontendParameters(request)
	assert.NotNil(t, err)
	assert.Equal(t, fv, &FrontendView{}, "structs should be equal + empty")

	// invalid Cookie
	cookieValue = "{'Month':10,'Year':2022,'TerminPerS,'MinDate':'2022-09-12T10:00:00Z'}"
	http.SetCookie(recorder, &http.Cookie{Name: "fe_parameter", Value: cookieValue})
	request = &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}
	assert.NotNil(t, err)
	assert.Equal(t, fv, &FrontendView{}, "structs should be equal + empty (invalid)")

}

func TestGetFeCookieString(t *testing.T) {
	fv := FrontendView{}
	str, err := GetFeCookieString(fv)
	assert.Nil(t, err, "Error should be nil")
	// only check for substring because one time.Now is executed later
	expSubstr := "{'Month':" + strconv.FormatInt(int64(time.Now().Month()), 10) + ",'Year':" + strconv.FormatInt(int64(time.Now().Year()), 10) + ",'TerminPerSite':7,'TerminSite':1,'MinDate':'"
	assert.True(t, strings.Contains(str, expSubstr), "should contain substring")
	fv = FrontendView{
		Month:         10,
		Year:          2022,
		TerminPerSite: 10,
		TerminSite:    2,
		MinDate:       time.Date(2022, 9, 12, 10, 00, 0, 0, time.UTC),
	}
	str, err = GetFeCookieString(fv)
	expStr := "{'Month':10,'Year':2022,'TerminPerSite':10,'TerminSite':2,'MinDate':'2022-09-12T10:00:00Z'}"
	assert.Equal(t, str, expStr, "strings should be equal")
	assert.Nil(t, err, "Error should be nil")
}

func TestGetTerminList(t *testing.T) {
	defer after()
	dataModel.InitDataModel("../data/test")
	user, _ := dataModel.Dm.AddUser("Testuser", "test", 0)
	ap1, ap2, ap3, ap4, ap5 := addAppointments(user.Id)
	// Termine ab 11.12.2022
	fv := FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2022, 11, 11, 11, 11, 1, 1, time.Local),
	}
	Apps := fv.GetTerminList(user.Appointments)
	// Richtige Reihenfolge der Termine
	exp := make([]data.Appointment, 0, 1)
	exp = append(exp, *ap3)
	exp = append(exp, *ap1)
	exp = append(exp, *ap2)
	expApp4 := *ap4
	expApp4.DateTimeStart = expApp4.DateTimeStart.AddDate(1, 0, 0)
	expApp4.DateTimeEnd = expApp4.DateTimeEnd.AddDate(1, 0, 0)
	exp = append(exp, expApp4)
	exp = append(exp, *ap5)
	assert.Equal(t, Apps, exp, "test 1 equal")

	// Termine ab 20.12.2022
	fv = FrontendView{
		Month:         12,
		Year:          2022,
		TerminSite:    1,
		TerminPerSite: 5,
		MinDate:       time.Date(2022, 12, 20, 11, 11, 1, 1, time.Local),
	}
	Apps = fv.GetTerminList(user.Appointments)
	// Richtige Reihenfolge der Termine
	exp = make([]data.Appointment, 0, 1)
	expApp5 := *ap5
	expApp5.DateTimeStart = expApp5.DateTimeStart.AddDate(0, 0, 3)
	expApp5.DateTimeEnd = expApp5.DateTimeEnd.AddDate(0, 0, 3)
	exp = append(exp, expApp5)
	expApp2 := *ap2
	expApp2.DateTimeStart = expApp2.DateTimeStart.AddDate(0, 0, 14)
	expApp2.DateTimeEnd = expApp2.DateTimeEnd.AddDate(0, 0, 14)
	exp = append(exp, expApp2)
	expApp3 := *ap3
	expApp3.DateTimeStart = expApp3.DateTimeStart.AddDate(0, 2, 0)
	expApp3.DateTimeEnd = expApp3.DateTimeEnd.AddDate(0, 2, 0)
	exp = append(exp, expApp3)
	expApp4 = *ap4
	expApp4.DateTimeStart = expApp4.DateTimeStart.AddDate(2, 0, 0)
	expApp4.DateTimeEnd = expApp4.DateTimeEnd.AddDate(2, 0, 0)
	exp = append(exp, expApp4)
	assert.Equal(t, Apps, exp, "test 2 equal")
}

func addAppointments(id int) (ap1, ap2, ap3, ap4, ap5 *data.Appointment) {
	// Einzelner Termin: 7 Dez 2022
	_, ap1 = dataModel.Dm.AddAppointment(id, "titel1", "beschreibung1", "here", time.Date(2022, 12, 7, 12, 0, 0, 0, time.Local),
		time.Date(2022, 12, 7, 14, 0, 0, 0, time.Local), false, 0, false)
	// Wöchentlicher Termin: ab 12 Dez 2022
	_, ap2 = dataModel.Dm.AddAppointment(id, "titel2", "beschreibung2", "here", time.Date(2022, 12, 12, 14, 0, 0, 0, time.Local),
		time.Date(2022, 12, 12, 17, 0, 0, 0, time.Local), true, 7, false)
	// Monatlicher Termin: ab 11 Nov 2022
	_, ap3 = dataModel.Dm.AddAppointment(id, "titel3", "beschreibung3", "here", time.Date(2022, 11, 11, 15, 0, 0, 0, time.Local),
		time.Date(2022, 11, 11, 17, 0, 0, 0, time.Local), true, 30, false)
	// Jährlicher Termin ab 12 Dez 2021
	_, ap4 = dataModel.Dm.AddAppointment(id, "titel4", "beschreibung4", "here", time.Date(2021, 12, 12, 15, 0, 0, 0, time.Local),
		time.Date(2021, 12, 12, 17, 0, 0, 0, time.Local), true, 365, false)
	// Täglicher Termin ab 17 Dez 2022
	_, ap5 = dataModel.Dm.AddAppointment(id, "titel5", "beschreibung5", "here", time.Date(2022, 12, 17, 15, 0, 0, 0, time.Local),
		time.Date(2022, 12, 17, 17, 0, 0, 0, time.Local), true, 1, false)
	return
}

func after() {
	os.RemoveAll("../data/test/")
	os.MkdirAll("../data/test/", 777)
}
