package error

import (
	"go_cal/templates"
	"net/http"
)

type ErrorType string

// add more error descriptions here
const (
	Default2          ErrorType = "Internal Server Error"
	Authentification  ErrorType = "Authentification failed"
	DuplicateUserName ErrorType = "Username already exists"
	WrongCredentials  ErrorType = "Username or password is wrong"
	InvalidInput      ErrorType = "Given input has wrong type"
	TitleIsEmpty      ErrorType = "Title of appointment is empty"
	EndBeforeBegin    ErrorType = "End date is earlier than start date"
)

type displayedError struct {
	Text string
	Link string
}

func CreateError(errorType ErrorType, prevLink string, w http.ResponseWriter, code int) {
	var error displayedError
	error = displayedError{
		Text: string(errorType),
		// TODO: http + host austauschen
		Link: "http://localhost:8080" + prevLink,
	}
	w.WriteHeader(code)
	templates.TempError.Execute(w, error)
	return
}
