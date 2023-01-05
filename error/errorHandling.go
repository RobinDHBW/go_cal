package error

type TypeError string

// add more error descriptions here
const (
	Default2          TypeError = "Internal Server Error"
	Authentication    TypeError = "Authentication failed"
	DuplicateUserName TypeError = "Username already exists"
	WrongCredentials  TypeError = "Username or password is wrong"
	InvalidInput      TypeError = "Given input has wrong type"
	TitleIsEmpty      TypeError = "Title of appointment is empty"
	EndBeforeBegin    TypeError = "End date is earlier than start date"
	EmptyField        TypeError = "Field for username/password is empty or usage of invalid characters (only alphanumeric and underscore are allowed)"
	InvalidUrl        TypeError = "Invalid Url"
)

// DisplayedError displayed Error in html
type DisplayedError struct {
	Text string
	Link string
}

// CreateError creates DisplayedError for displaying to user
func CreateError(errorType TypeError, prevLink string) (error DisplayedError) {
	error = DisplayedError{
		Text: string(errorType),
		Link: prevLink,
	}
	return error
}
