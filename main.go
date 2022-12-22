package main

import (
	"html/template"
	"log"
	"net/http"
)

type data struct {
	Name  string
	Email string
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")

	d := data{
		Name:  name,
		Email: "noreply@xyz.de",
	}
	var tempInit = template.Must(template.ParseFiles("./templates/test.html"))
	tempInit.Execute(w, d)
}
func main() {
	http.HandleFunc("/", mainHandler)
	http.Handle("/templates/static/", http.StripPrefix("/templates/static", http.FileServer(http.Dir("templates/static"))))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
