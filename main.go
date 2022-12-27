package main

import (
	"fmt"
	"go_cal/authentication"
	"go_cal/calendarView"
	"go_cal/templates"
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	templates.Init()
	err := authentication.LoadUsersFromFiles()
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/updateCalendar", calendarView.UpdateCalendarHandler)
	http.HandleFunc("/register", authentication.RegisterHandler)
	http.HandleFunc("/logout", authentication.LogoutHandler)
	http.HandleFunc("/", authentication.LoginHandler)

	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
