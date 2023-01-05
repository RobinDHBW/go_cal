package templates

import (
	"html/template"
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
var TempTerminVotingSuccess *template.Template
var TempSearchTermin *template.Template

// Init initialize html templates
// dir has to be the root path
func Init(dir string) {
	TempInit = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/calendar.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempLogin = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/login.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempRegister = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/register.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempError = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/error.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempTerminList = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminlist.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempTerminEdit = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminedit.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempCreateTermin = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/termincreate.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempShareTermin = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminshare.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempCreateShareTermin = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminsharecreate.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempEditShareTermin = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminshareedit.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempSearchTermin = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminsearch.tmpl.html"), filepath.Join(dir, "/templates/header.tmpl.html")))
	TempTerminVoting = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminvoting.tmpl.html")))
	TempTerminVotingSuccess = template.Must(template.ParseFiles(filepath.Join(dir, "/templates/terminvotingsuccess.tmpl.html")))
}
