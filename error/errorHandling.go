package error

type ErrorType string

// add more error descriptions here
const (
	Default2          ErrorType = "Internal Server Error"
	Authentification  ErrorType = "Authentification failed"
	DuplicateUserName ErrorType = "Username already exists"
	WrongCredentials  ErrorType = "Username or password is wrong"
)

type DisplayedError struct {
	Text string
	Link string
}

func CreateError(errorType ErrorType, prevLink string) (error DisplayedError) {
	error = DisplayedError{
		Text: string(errorType),
		// TODO: http + host austauschen
		Link: "http://localhost:8080" + prevLink,
	}
	return error
}
