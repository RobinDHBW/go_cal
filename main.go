package main

import (
	"go_cal/authentication"
	"go_cal/calendar"
	"go_cal/configuration"
	"go_cal/dataModel"
	"go_cal/templates"
	"go_cal/terminHandling"
	"log"
	"net/http"
	"strconv"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// Flags einlesen
	configuration.ReadFlags()
	// Datamodel initialisieren
	dataModel.InitDataModel(configuration.Folder)
	// Server für Channel-Kommunikation initialisieren
	authentication.InitServer()
	// html templates initialisieren
	templates.Init()
	// setzt rand.Seed für Generierung von einmaligen Tokens
	dataModel.InitSeed()

	// Endpunkte definieren
	http.HandleFunc("/updateCalendar", authentication.Wrapper(calendar.UpdateCalendarHandler))
	http.HandleFunc("/register", authentication.RegisterHandler)
	http.HandleFunc("/logout", authentication.LogoutHandler)
	http.HandleFunc("/", authentication.LoginHandler)
	// TODO Wrapper aufrufen
	http.HandleFunc("/listTermin", authentication.Wrapper(terminHandling.TerminHandler))
	http.HandleFunc("/createTermin", authentication.Wrapper(terminHandling.TerminCreateHandler))
	http.HandleFunc("/editTermin", authentication.Wrapper(terminHandling.TerminEditHandler))
	http.HandleFunc("/listShareTermin", authentication.Wrapper(terminHandling.TerminShareListHandler))
	http.HandleFunc("/shareTermin", authentication.Wrapper(terminHandling.TerminShareHandler))
	http.HandleFunc("/terminVoting", terminHandling.TerminVotingHandler)
	//http.HandleFunc("/download", export.Wrapper(export.AuthenticatorFunc(export.CheckUserValid), terminHandling.DownloadHandler))
	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))

	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(configuration.Port), nil))
}
