package templates

import (
	"html/template"
	"os"
	"path/filepath"
)

var TempInit *template.Template
var TempLogin *template.Template
var TempRegister *template.Template
var TempError *template.Template
var TempTerminList *template.Template
var TempTerminEdit *template.Template
var TempCreateTermin *template.Template
var TempShareTermin *template.Template
var TempCreateShareTermin *template.Template
var TempEditShareTermin *template.Template
var TempTerminVoting *template.Template

func Init() {
	dir, _ := os.Getwd()
	if filepath.Base(dir) == "go_cal" || filepath.Base(dir) == "Go-Kalender" {
		TempInit = template.Must(template.ParseFiles("./templates/calendar.tmpl.html", "./templates/header.tmpl.html"))
		TempLogin = template.Must(template.ParseFiles("./templates/login.tmpl.html", "./templates/header.tmpl.html"))
		TempRegister = template.Must(template.ParseFiles("./templates/register.tmpl.html", "./templates/header.tmpl.html"))
		TempError = template.Must(template.ParseFiles("./templates/error.tmpl.html", "./templates/header.tmpl.html"))
		TempTerminList = template.Must(template.ParseFiles("./templates/terminlist.tmpl.html", "./templates/header.tmpl.html"))
		TempTerminEdit = template.Must(template.ParseFiles("./templates/terminedit.tmpl.html", "./templates/header.tmpl.html"))
		TempCreateTermin = template.Must(template.ParseFiles("./templates/termincreate.tmpl.html", "./templates/header.tmpl.html"))
		TempShareTermin = template.Must(template.ParseFiles("./templates/terminshare.tmpl.html", "./templates/header.tmpl.html"))
		TempCreateShareTermin = template.Must(template.ParseFiles("./templates/terminsharecreate.tmpl.html", "./templates/header.tmpl.html"))
		TempEditShareTermin = template.Must(template.ParseFiles("./templates/terminshareedit.tmpl.html", "./templates/header.tmpl.html"))
		TempTerminVoting = template.Must(template.ParseFiles("./templates/terminvoting.tmpl.html"))
	} else {
		TempInit = template.Must(template.ParseFiles("../templates/calendar.tmpl.html", "../templates/header.tmpl.html"))
		TempLogin = template.Must(template.ParseFiles("../templates/login.tmpl.html", "../templates/header.tmpl.html"))
		TempRegister = template.Must(template.ParseFiles("../templates/register.tmpl.html", "../templates/header.tmpl.html"))
		TempError = template.Must(template.ParseFiles("../templates/error.tmpl.html", "../templates/header.tmpl.html"))
		TempTerminList = template.Must(template.ParseFiles("../templates/terminlist.tmpl.html", "../templates/header.tmpl.html"))
		TempTerminEdit = template.Must(template.ParseFiles("../templates/terminedit.tmpl.html", "../templates/header.tmpl.html"))
		TempCreateTermin = template.Must(template.ParseFiles("../templates/termincreate.tmpl.html", "../templates/header.tmpl.html"))
		TempShareTermin = template.Must(template.ParseFiles("../templates/terminshare.tmpl.html", "../templates/header.tmpl.html"))
		TempCreateShareTermin = template.Must(template.ParseFiles("../templates/terminsharecreate.tmpl.html", "../templates/header.tmpl.html"))
		TempEditShareTermin = template.Must(template.ParseFiles("../templates/terminshareedit.tmpl.html", "../templates/header.tmpl.html"))
		TempTerminVoting = template.Must(template.ParseFiles("../templates/terminvoting.tmpl.html"))
	}
}
