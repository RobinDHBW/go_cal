package templates

import "html/template"

var TempInit = template.Must(template.ParseFiles("./templates/test.tmpl.html"))
var TempLogin = template.Must(template.ParseFiles("./templates/login.tmpl.html"))
var TempRegister = template.Must(template.ParseFiles("./templates/register.tmpl.html"))
var TempError = template.Must(template.ParseFiles("./templates/error.tmpl.html"))
