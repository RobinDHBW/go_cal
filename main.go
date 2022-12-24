package main

import (
	"fmt"
	"go_cal/authentication"
	"go_cal/calendarView"
	"go_cal/templates"
	"log"
	"net/http"
)

// initialize ErrorList
var errorList = make(map[string]string)

type displayedError struct {
	Text string
	Link string
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	configErrorList()
	err := authentication.LoadUsersFromFiles()
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/updateCalendar", calendarView.UpdateCalendarHandler)
	http.HandleFunc("/register", authentication.RegisterHandler)
	http.HandleFunc("/error", ErrorHandler)

	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))
	http.HandleFunc("/", authentication.LoginHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func configErrorList() {
	errorList["default"] = "Internal Server Error"
	errorList["authentification"] = "Authentification failed"
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	var error displayedError
	errorType := r.URL.Query().Get("type")
	prevLink := r.URL.Query().Get("link")

	value, ok := errorList[errorType]
	if ok {
		error = displayedError{
			Text: value,
			// TODO: http austauschen
			Link: "http://" + r.Host + prevLink,
		}
	} else {
		error = displayedError{
			Text: errorList["default"],
			// TODO: http austauschen
			Link: "http://" + r.Host + "/",
		}
	}
	templates.TempError.Execute(w, error)
	return
}
