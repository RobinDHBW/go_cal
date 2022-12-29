package authentication

// session authentification inspired from https://github.com/sohamkamani/go-session-auth-example
// channel communication inspried from https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/
// https://github.com/eliben/code-for-blog/blob/master/2019/gohttpconcurrency/channel-manager-server.go

import (
	"encoding/json"
	"fmt"
	error2 "go_cal/error"
	"go_cal/templates"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

// TODO: has to be removed, use datamodel
// map with username and corresponding hashed password
var users = map[string][]byte{}

// TODO: has to be removed, use datamodel
// Credentials struct for a user
type Credentials struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

// session consist of n user and an expiry time
type session struct {
	uname   string
	expires time.Time
}

//// map with SessionTokens and corresponding sessions
//var sessions = map[string]*session{}

// prüft ob Session abgelaufen ist
func (s session) isExpired() bool {
	return s.expires.Before(time.Now())
}

type Server struct {
	Cmds chan<- Command
}

type CommandType string

const (
	read   CommandType = "read"
	write  CommandType = "write"
	remove CommandType = "remove"
	update CommandType = "update"
)

type Command struct {
	ty           CommandType
	sessionToken string
	session      *session
	replyChannel chan *session
}

func StartSessionManager() chan<- Command {
	// map with SessionTokens and corresponding sessions
	sessions := map[string]*session{}

	cmds := make(chan Command)

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
				sessions[cmd.sessionToken].expires = time.Now().Add(1 * time.Minute)
				cmd.replyChannel <- sessions[cmd.sessionToken]
			}
		}
	}()
	return cmds
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie überprüfen
	isCookieValid := checkCookie(r, s)

	// kein gültiger Cookie im Request --> login-procedure
	if !isCookieValid {
		// übermitteltes Formular parsen
		r.ParseForm()
		// wenn Login-Button gedrückt und POST ausgeführt wurde
		if r.PostForm.Has("login") && r.Method == http.MethodPost {
			// Eingabefelder (username und password) auslesen
			password := []byte(r.PostFormValue("passwd"))
			username := r.PostFormValue("uname")
			// Eingabevalidierung
			valid := validateInput(username, password)
			if !valid {
				// Response header 400 setzen
				w.WriteHeader(http.StatusBadRequest)
				// Fehlermeldung für Nutzer anzeigen
				templates.TempError.Execute(w, error2.CreateError(error2.EmptyField, r.Host+"/"))
				return
			}
			// user authentifizieren
			successful := AuthenticateUser(username, password)
			// user erfolgreich authentifiziert
			if successful {
				// neue session erstellen
				sessionToken, expires := createSession(username, s)
				// Cookie in response setzen
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: expires,
				})
				// redirect auf Kalender
				http.Redirect(w, r, "/updateCalendar", http.StatusFound)
				return
				// user nicht erfolgreich authentifiziert (username oder password falsch)
			} else {
				// Response header 401 setzen
				w.WriteHeader(http.StatusUnauthorized)
				// Fehlermeldung für Nutzer anzeigen
				templates.TempError.Execute(w, error2.CreateError(error2.WrongCredentials, r.Host+"/"))
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

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// übermitteltes Formular parsen
	r.ParseForm()
	// wenn Register-Button gedrückt und POST ausgeführt wurde
	if r.PostForm.Has("register") && r.Method == http.MethodPost {
		// Eingabefelder (username und password) auslesen
		password := []byte(r.PostFormValue("passwd"))
		username := r.PostFormValue("uname")
		// Eingabevalidierung
		valid := validateInput(username, password)
		if !valid {
			// Response header 400 setzen
			w.WriteHeader(http.StatusBadRequest)
			// Fehlermeldung für Nutzer anzeigen
			templates.TempError.Execute(w, error2.CreateError(error2.EmptyField, r.Host+"/register"))
			return
		}
		// exisitiert der Nutzername schon?
		duplicate := isDuplicateUsername(username)
		// wenn Nutzername schon exisitert
		if duplicate {
			// Response header 400 setzen
			w.WriteHeader(http.StatusBadRequest)
			// Fehlermeldung für Nutzer anzeigen
			templates.TempError.Execute(w, error2.CreateError(error2.DuplicateUserName, r.Host+"/register"))
			return
			// Nutzername exisitert noch nicht --> register möglich
		} else {
			// Passwort hashen
			hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
			user := Credentials{
				Username: username,
				Password: hashedPassword,
			}

			// ab hier Speicherprozess
			folder, err := os.Open("./files")
			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = folder.Readdir(0)
			if err != nil {
				fmt.Println(err)
				return
			}

			text, _ := json.Marshal(user)

			file, err := os.Create("./files/" + user.Username + ".json")
			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = file.Write(text)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = file.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			// neue session erstellen
			sessionToken, expires := createSession(username, s)
			// Cookie in response setzen
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expires,
			})
			// sync file to variables
			LoadUsersFromFiles()
			// redirect auf Kalender
			http.Redirect(w, r, "/updateCalendar", http.StatusFound)
			return
		}
	}
	// Register-Seite ausliefern
	templates.TempRegister.Execute(w, nil)
}

// Wrapper für Authentifizierung mit Cookie
func (s *Server) Wrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Cookie aus Request überprüfen
		isCookieValid := checkCookie(r, s)
		// wenn Cookie valid
		if isCookieValid {
			// Cookie verlängern
			sessionToken, expires := refreshCookie(r, s)
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
			templates.TempError.Execute(w, error2.CreateError(error2.Authentification, r.Host+"/"))
			return
		}
	}
}

func refreshCookie(r *http.Request, s *Server) (sessionToken string, expires time.Time) {
	replyChannel := make(chan *session)
	// Cookie auslesen
	cookie, _ := r.Cookie("session_token")
	// Sessiontoken auslesen
	sessionToken = cookie.Value
	s.Cmds <- Command{ty: update, sessionToken: sessionToken, replyChannel: replyChannel}
	session := <-replyChannel
	// session auslesen
	//session, _ := sessions[sessionToken]
	// Session ist valide, da zuvor CheckCookie ausgeführt wurde
	// expires um 10 min verlägern
	//session.expires = session.expires.Add(1 * time.Minute)
	return sessionToken, session.expires
}

func createUUID(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func LoadUsersFromFiles() error {
	// open folder
	folder, err := os.Open("./files")
	if err != nil {
		return err
	}
	// read all files inside directory
	files, err := folder.Readdir(0)
	if err != nil {
		return err
	}
	for _, file := range files {
		var user Credentials
		data, _ := os.ReadFile("./files/" + file.Name())
		err := json.Unmarshal(data, &user)
		if err != nil {
			return err
		}
		users[user.Username] = user.Password
	}
	return nil
}

func AuthenticateUser(username string, unHashedPassword []byte) (successful bool) {
	expectedPasswordHash, ok := users[username]
	if !ok || bcrypt.CompareHashAndPassword(expectedPasswordHash, unHashedPassword) != nil {
		return false
	} else {
		return true
	}
}

func checkCookie(r *http.Request, s *Server) (successful bool) {
	replyChannel := make(chan *session)

	// Cookie auslesen
	cookie, err := r.Cookie("session_token")
	// kein Cookie
	if err == http.ErrNoCookie {
		return false
	}
	// Sessiontoken auslesen
	sessionToken := cookie.Value
	s.Cmds <- Command{ty: read, sessionToken: sessionToken, replyChannel: replyChannel}
	session := <-replyChannel

	//if session == nil {
	//	return false
	//}

	//// session auslesen
	//session, ok := sessions[sessionToken]
	//// keine Session zu Sessiontoken gefunden
	//if !ok {
	//	return false
	//}

	// SessionToken is abgelaufen
	if session.isExpired() {
		// Session löschen
		s.Cmds <- Command{ty: remove, sessionToken: sessionToken, replyChannel: replyChannel}
		//delete(sessions, sessionToken)
		<-replyChannel
		return false
	}
	return true
}

func isDuplicateUsername(username string) (isDuplicate bool) {
	// existiert der username schon?
	_, ok := users[username]
	return ok
}

func createSession(username string, s *Server) (sessionToken string, expires time.Time) {
	// Sessiontoken generieren
	sessionToken = createUUID(25)
	// Session läuft nach x Minuten ab
	// TODO Zeit anpassen
	expires = time.Now().Add(1 * time.Minute)
	// Session anhand des Sessiontokens speichern
	replyChannel := make(chan *session)
	s.Cmds <- Command{ty: write, sessionToken: sessionToken, session: &session{uname: username, expires: expires}, replyChannel: replyChannel}
	//sessions[sessionToken] = &session{
	//	uname:   username,
	//	expires: expires,
	//}
	session := <-replyChannel
	return sessionToken, session.expires
}

// Überprüft Nutzereingaben beim Login und Registrieren
func validateInput(username string, password []byte) (successful bool) {
	// wenn Felder leer
	if len(username) == 0 || len(password) == 0 {
		return false
	}
	// wenn unerlaubte Zeichen verwendet werden
	const invalidCharactersUsername string = "[\\\\/:*?\"<>|{}`´']"
	//const invalidCharactersPassword string = "[<>{}`´']"
	matchUsername, _ := regexp.MatchString(invalidCharactersUsername, username)
	matchPassword, _ := regexp.MatchString(invalidCharactersUsername, string(password))
	if matchUsername || matchPassword {
		return false
	}
	return true
}
