package main

import (
	"go_cal/calendarView"
	"html/template"
	"log"
	"net/http"
	"time"
)

type data struct {
	Name  string
	Email string
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	//q := r.URL.Query()
	//name := q.Get("name")

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
	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
