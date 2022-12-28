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

func Init() {
	dir, _ := os.Getwd()
	if filepath.Base(dir) == "go_cal" {
		TempInit = template.Must(template.ParseFiles("./templates/test.tmpl.html", "./templates/header.tmpl.html"))
		TempLogin = template.Must(template.ParseFiles("./templates/login.tmpl.html", "./templates/header.tmpl.html"))
		TempRegister = template.Must(template.ParseFiles("./templates/register.tmpl.html", "./templates/header.tmpl.html"))
		TempError = template.Must(template.ParseFiles("./templates/error.tmpl.html", "./templates/header.tmpl.html"))
		TempTerminList = template.Must(template.ParseFiles("./templates/terminlist.tmpl.html", "./templates/header.tmpl.html"))
		TempTerminEdit = template.Must(template.ParseFiles("./templates/terminedit.tmpl.html", "./templates/header.tmpl.html"))
	} else {
		TempInit = template.Must(template.ParseFiles("../templates/test.tmpl.html", "../templates/header.tmpl.html"))
		TempLogin = template.Must(template.ParseFiles("../templates/login.tmpl.html", "../templates/header.tmpl.html"))
		TempRegister = template.Must(template.ParseFiles("../templates/register.tmpl.html", "../templates/header.tmpl.html"))
		TempError = template.Must(template.ParseFiles("../templates/error.tmpl.html", "../templates/header.tmpl.html"))
		TempTerminList = template.Must(template.ParseFiles("../templates/terminlist.tmpl.html", "../templates/header.tmpl.html"))
		TempTerminEdit = template.Must(template.ParseFiles("../templates/terminedit.tmpl.html", "../templates/header.tmpl.html"))
	}
}
