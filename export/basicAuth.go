// Matrikelnummern:
// 9495107, 4706893, 9608900

package export

import (
	"go_cal/authentication"
	"net/http"
)

// an Vorlesung orientiert

type Authenticator interface {
	Authenticate(user, password string) bool
}

type AuthenticatorFunc func(user, password string) bool

func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

// CheckUserValid checks whether the username with the given password is authorized
func CheckUserValid(username, password string) bool {
	successful := authentication.AuthenticateUser(username, password)
	return successful
}

// Wrapper for Basic Auth
func Wrapper(authenticator Authenticator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, passwd, ok := r.BasicAuth()
		if ok && authenticator.Authenticate(user, passwd) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate",
				"Basic realm=\"My Simple Server\"")
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		}
	}
}
