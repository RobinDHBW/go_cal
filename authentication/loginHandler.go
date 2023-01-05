// Matrikelnummern:
// 9495107, 4706893, 9608900

package authentication

// session authentication inspired from https://github.com/sohamkamani/go-session-auth-example
// channel communication inspried from https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/

import (
	"errors"
	"go_cal/configuration"
	"go_cal/data"
	"go_cal/dataModel"
	error2 "go_cal/error"
	"go_cal/frontendHandling"
	"go_cal/templates"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

// Server for channel communication
type Server struct {
	Cmds chan<- Command
}

var Serv *Server

// session consists of username and expiry time
type session struct {
	uname   string
	expires time.Time
}

// CommandType for channel communication (read, write, remove, update)
type CommandType string

const (
	read   CommandType = "read"
	write  CommandType = "write"
	remove CommandType = "remove"
	update CommandType = "update"
)

// Command for channel communication
// conatains necessary information
type Command struct {
	ty           CommandType
	sessionToken string
	session      *session
	replyChannel chan *session
}

// InitServer sets global var Serv
// execute StartSessionManager
func InitServer() {
	Serv = &Server{
		Cmds: StartSessionManager(),
	}
}

// isExpired checks whether a session is expired
func (s session) isExpired() bool {
	return s.expires.Before(time.Now())
}

// StartSessionManager starts a go routine to handle parallel access to the session map
// communication is done via channels
// read, write, remove or update a session
func StartSessionManager() chan<- Command {
	// map with SessionTokens and corresponding sessions, cannot be accessed from outside
	sessions := map[string]*session{}
	cmds := make(chan Command)
	// start go routine
	// listen on incoming Commands (read, write, remove, update)
	// write result in replyChannel
	go func() {
		for cmd := range cmds {
			switch cmd.ty {
			case read:
				if val, ok := sessions[cmd.sessionToken]; ok {
					cmd.replyChannel <- val
				} else {
					cmd.replyChannel <- &session{}
				}
			case write:
				sessions[cmd.sessionToken] = cmd.session
				cmd.replyChannel <- cmd.session
			case remove:
				delete(sessions, cmd.sessionToken)
				cmd.replyChannel <- &session{}
			case update:
				sessions[cmd.sessionToken].expires = time.Now().Add(time.Minute * time.Duration(configuration.Timeout))
				cmd.replyChannel <- sessions[cmd.sessionToken]
			}
		}
	}()
	return cmds
}

// LoginHandler handle inputs of login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie überprüfen
	isCookieValid := checkCookie(r)

	// kein gültiger Cookie im Request --> login-procedure
	if !isCookieValid {
		// übermitteltes Formular parsen
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/listTermin"))
			return
		}
		// wenn Login-Button gedrückt und POST ausgeführt wurde
		if r.PostForm.Has("login") && r.Method == http.MethodPost {
			// Eingabefelder (username und password) auslesen
			password := r.PostFormValue("passwd")
			username := r.PostFormValue("uname")
			// Eingabevalidierung
			valid := validateInput(username, password)
			if !valid {
				// Response header 400 setzen
				w.WriteHeader(http.StatusBadRequest)
				// Fehlermeldung für Nutzer anzeigen
				templates.TempError.Execute(w, error2.CreateError(error2.EmptyField, "/"))
				return
			}
			// user authentifizieren
			successful := AuthenticateUser(username, password)
			// user erfolgreich authentifiziert
			if successful {
				// neue session erstellen
				sessionToken, expires := CreateSession(username)
				// Cookie in response setzen
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: expires,
				})
				cookieValue, err := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/"))
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:  "fe_parameter",
					Value: cookieValue,
				})
				// redirect auf Kalender
				http.Redirect(w, r, "/updateCalendar", http.StatusFound)
				return
				// user nicht erfolgreich authentifiziert (username oder password falsch)
			} else {
				// Response header 401 setzen
				w.WriteHeader(http.StatusUnauthorized)
				// Fehlermeldung für Nutzer anzeigen
				templates.TempError.Execute(w, error2.CreateError(error2.WrongCredentials, "/"))
				return
			}
		}
		// gültiger Cookie im Request --> kein Login nötig
	} else {
		// redirect auf Kalender
		http.Redirect(w, r, "/updateCalendar", http.StatusFound)
		return
	}
	// Login-Seite ausliefern
	templates.TempLogin.Execute(w, nil)
}

// RegisterHandler handle inputs of registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// übermitteltes Formular parsen
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.TempError.Execute(w, error2.CreateError(error2.Default2, "/listTermin"))
		return
	}
	// wenn Register-Button gedrückt und POST ausgeführt wurde
	if r.PostForm.Has("register") && r.Method == http.MethodPost {
		// Eingabefelder (username und password) auslesen
		password := r.PostFormValue("passwd")
		username := r.PostFormValue("uname")
		// Eingabevalidierung
		valid := validateInput(username, password)
		if !valid {
			// Response header 400 setzen
			w.WriteHeader(http.StatusBadRequest)
			// Fehlermeldung für Nutzer anzeigen
			templates.TempError.Execute(w, error2.CreateError(error2.EmptyField, "/register"))
			return
		}
		// neuen User erstellen
		_, err := dataModel.Dm.AddUser(username, password, 1)
		// Nutzername existiert schon, Erstellung war nicht erfolgreich
		if err != nil {
			// Response header 400 setzen
			w.WriteHeader(http.StatusBadRequest)
			// Fehlermeldung für Nutzer anzeigen
			templates.TempError.Execute(w, error2.CreateError(error2.DuplicateUserName, "/register"))
			return
			// Nutzername existiert noch nicht, Erstellung war erfolgreich
		} else {
			// neue session erstellen
			sessionToken, expires := CreateSession(username)
			// Cookie in response setzen
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expires,
			})
			cookieValue, err := frontendHandling.GetFeCookieString(frontendHandling.FrontendView{})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				templates.TempError.Execute(w, error2.CreateError(error2.InvalidInput, "/"))
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "fe_parameter",
				Value: cookieValue,
			})
			// redirect auf Kalender
			http.Redirect(w, r, "/updateCalendar", http.StatusFound)
			return
		}
	}
	// Register-Seite ausliefern
	templates.TempRegister.Execute(w, nil)
}

// Wrapper für Authentifizierung mit Cookie
func Wrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Cookie aus Request überprüfen
		isCookieValid := checkCookie(r)
		// wenn Cookie valid
		if isCookieValid {
			// Cookie verlängern
			sessionToken, expires := refreshCookie(r)
			// Cookie setzen
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expires,
			})
			// handler aufrufen
			handler(w, r)
			// wenn Cookie invalide
		} else {
			// Response header 401 setzen
			w.WriteHeader(http.StatusUnauthorized)
			// Fehlermeldung für Nutzer anzeigen
			templates.TempError.Execute(w, error2.CreateError(error2.Authentication, "/"))
			return
		}
	}
}

// refreshCookie refreshes provided cookie in request
// increase expires
func refreshCookie(r *http.Request) (sessionToken string, expires time.Time) {
	replyChannel := make(chan *session)
	// Cookie auslesen
	cookie, _ := r.Cookie("session_token")
	// Sessiontoken auslesen
	sessionToken = cookie.Value
	// Session aktualisieren
	Serv.Cmds <- Command{ty: update, sessionToken: sessionToken, replyChannel: replyChannel}
	replySession := <-replyChannel
	return sessionToken, replySession.expires
}

// createUUID creates alphabetical ID for sessionToken
func createUUID(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// AuthenticateUser checks whether the provided username and password are correct
func AuthenticateUser(username, unHashedPassword string) (successful bool) {
	user := dataModel.Dm.GetUserByName(username)
	if user != nil && dataModel.Dm.ComparePW(unHashedPassword, user.Password) {
		return true
	} else {
		return false
	}
}

// checkCookie checks whether the provided cookie in the request is valid
func checkCookie(r *http.Request) (successful bool) {
	// Anwortchannel erstellen
	replyChannel := make(chan *session)
	// Cookie auslesen
	cookie, err := r.Cookie("session_token")
	// kein Cookie
	if err == http.ErrNoCookie {
		return false
	}
	// Sessiontoken auslesen
	sessionToken := cookie.Value
	// read-Command schicken
	Serv.Cmds <- Command{ty: read, sessionToken: sessionToken, replyChannel: replyChannel}
	// session aus Antwortchannel lesen
	replySession := <-replyChannel
	// SessionToken is abgelaufen
	if replySession.isExpired() {
		// Session löschen
		Serv.Cmds <- Command{ty: remove, sessionToken: sessionToken, replyChannel: replyChannel}
		<-replyChannel
		return false
	}
	return true
}

// CreateSession creates a session for a username
func CreateSession(username string) (sessionToken string, expires time.Time) {
	// Anwortchannel erstellen
	replyChannel := make(chan *session)
	// Sessiontoken generieren
	sessionToken = createUUID(25)
	// Session läuft nach x Minuten ab (über flag gesteuert)
	expires = time.Now().Add(time.Minute * time.Duration(configuration.Timeout))
	// Session anhand des Sessiontokens speichern
	Serv.Cmds <- Command{ty: write, sessionToken: sessionToken, session: &session{uname: username, expires: expires}, replyChannel: replyChannel}
	// session aus Antwortchannel lesen
	replySession := <-replyChannel
	return sessionToken, replySession.expires
}

// validateInput Checks user input during login and registration
func validateInput(username, password string) (successful bool) {
	// wenn Felder leer
	if len(username) == 0 || len(password) == 0 {
		return false
	}
	// wenn unerlaubte Zeichen verwendet werden
	const validCharacters string = "^[a-zA-Z0-9_]*$"
	matchUsername, _ := regexp.MatchString(validCharacters, username)
	matchPassword, _ := regexp.MatchString(validCharacters, password)
	if !matchUsername || !matchPassword {
		return false
	}
	return true
}

// GetUserBySessionToken returns user from cookie inside a request
func GetUserBySessionToken(r *http.Request) (*data.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	sessionToken := cookie.Value
	replyChannel := make(chan *session)
	Serv.Cmds <- Command{ty: read, sessionToken: sessionToken, replyChannel: replyChannel}
	replySession := <-replyChannel
	if *replySession == (session{}) {
		return nil, errors.New("cannot get User")
	}
	username := replySession.uname
	user := dataModel.Dm.GetUserByName(username)
	if user == nil {
		return nil, errors.New("cannot get User")
	}
	return user, nil
}

// getUsernameBySessionToken returns username from corresponding session token
func getUsernameBySessionToken(sessionToken string) string {
	replyChannel := make(chan *session)
	Serv.Cmds <- Command{ty: read, sessionToken: sessionToken, replyChannel: replyChannel}
	replySession := <-replyChannel
	if *replySession == (session{}) {
		return ""
	}
	return replySession.uname
}
