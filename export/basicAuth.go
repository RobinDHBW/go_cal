package export

import "net/http"

// an Vorlseung orientiert

type Authenticator interface {
	Authenticate(user, password string) bool
}

type AuthenticatorFunc func(user, password string) bool

func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

func checkUserValid(username, password string) bool {
	return true
}

func Wrapper(authenticator Authenticator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pswd, ok := r.BasicAuth()
		if ok && checkUserValid(user, pswd) {
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
