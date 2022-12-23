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

func mainHandler(w http.ResponseWriter, r *http.Request) {
	cal := calendarView.Calendar{
		Month:   time.Now().Month(),
		Year:    time.Now().Year(),
		Current: time.Now(),
	}
	var tempInit = template.Must(template.ParseFiles("./templates/test.tmpl.html"))
	tempInit.Execute(w, cal)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/updateCalendar", calendarView.UpdateCalendarHandler)
	http.HandleFunc("/login", authentication.LoginHandler)
	http.HandleFunc("/register", authentication.RegisterHandler)
	http.HandleFunc("/terminlist", terminHandling.TerminHandler)
	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
