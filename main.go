package main

import (
	"go_cal/authentication"
	"go_cal/calendarView"
	"go_cal/terminHandling"
	"html/template"
	"log"
	"net/http"
	"time"
)

var globalTemp = 0

func mainHandler(w http.ResponseWriter, r *http.Request) {
	cal := calendarView.Calendar{
		Month:   time.Now().Month(),
		Year:    time.Now().Year(),
		Current: time.Now(),
	}
	if globalTemp == 0 { // nur zum Testen
		terminHandling.TView.TList.CreateTermin("T1", "1 content", time.Now().AddDate(0, 0, -1), time.Now())
		terminHandling.TView.TList.CreateTermin("T2", "2 content", time.Now(), time.Now())
		terminHandling.TView.TList.CreateTermin("T3", "3 content", time.Now().AddDate(0, 0, -2), time.Now())
	}
	globalTemp = 1

	var tempInit = template.Must(template.ParseFiles("./templates/test.tmpl.html"))
	tempInit.Execute(w, cal)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/updateCalendar", calendarView.UpdateCalendarHandler)
	http.HandleFunc("/login", authentication.LoginHandler)
	http.HandleFunc("/register", authentication.RegisterHandler)
	http.HandleFunc("/terminlist", terminHandling.TerminHandler)
	http.HandleFunc("/updateTerminList", terminHandling.TerminHandler)
	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
