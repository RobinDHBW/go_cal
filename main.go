package main

import (
	"go_cal/authentication"
	"go_cal/calendarView"
	"go_cal/dataModel"
	"go_cal/templates"
	"go_cal/terminHandling"
	"log"
	"net/http"
	"time"
)

var globalTemp = 0

func mainHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	dataModel.InitDataModel("./files")
	authentication.InitServer()
	//authentication.Serv = &authentication.Server{Cmds: authentication.StartSessionManager()}
	templates.Init()
	if globalTemp == 0 { // nur zum Testen
		terminHandling.TView.TList.CreateTermin("T1", "1 content", time.Now().AddDate(0, 0, -1), time.Now(), terminHandling.None)
		terminHandling.TView.TList.CreateTermin("T2", "2 content", time.Now(), time.Now(), terminHandling.None)
		terminHandling.TView.TList.CreateTermin("T3", "3 content", time.Now().AddDate(0, 0, -2), time.Now(), terminHandling.None)
	}
	globalTemp = 1

	http.HandleFunc("/updateCalendar", authentication.Wrapper(calendarView.UpdateCalendarHandler))
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
