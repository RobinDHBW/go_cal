package main

import (
	"go_cal/authentication"
	"go_cal/calendar"
	"go_cal/dataModel"
	"go_cal/templates"
	"go_cal/terminHandling"
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	dataModel.InitDataModel("./files")
	authentication.InitServer()
	//authentication.Serv = &authentication.Server{Cmds: authentication.StartSessionManager()}
	templates.Init()

	http.HandleFunc("/updateCalendar", authentication.Wrapper(calendar.UpdateCalendarHandler))
	http.HandleFunc("/register", authentication.RegisterHandler)
	http.HandleFunc("/logout", authentication.LogoutHandler)
	http.HandleFunc("/", authentication.LoginHandler)
	http.HandleFunc("/listTermin", terminHandling.TerminHandler)
	http.HandleFunc("/createTermin", terminHandling.TerminCreateHandler)
	http.HandleFunc("/editTermin", terminHandling.TerminEditHandler)
	//http.HandleFunc("/download", export.Wrapper(export.AuthenticatorFunc(export.CheckUserValid), terminHandling.DownloadHandler))
	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
