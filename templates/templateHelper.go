package templates

import "html/template"

var TempInit = template.Must(template.ParseFiles("./templates/test.tmpl.html"))
var TempLogin = template.Must(template.ParseFiles("./templates/login.tmpl.html"))
var TempRegister = template.Must(template.ParseFiles("./templates/register.tmpl.html"))
var TempTerminList = template.Must(template.ParseFiles("./templates/terminlist.tmpl.html"))
var TempTerminEdit = template.Must(template.ParseFiles("./templates/terminedit.tmpl.html"))
