package authentication

// session authentification inspired from https://github.com/sohamkamani/go-session-auth-example

import (
	"encoding/json"
	"fmt"
	error2 "go_cal/error"
	"go_cal/templates"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"os"
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

// map with SessionTokens and corresponding sessions
var sessions = map[string]*session{}

// prüft ob Session abgelaufen ist
func (s session) isExpired() bool {
	return s.expires.Before(time.Now())
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie überprüfen
	isCookieValid := checkCookie(r)

	// kein gültiger Cookie im Request --> login-procedure
	if !isCookieValid {
		// übermitteltes Formular parsen
		r.ParseForm()
		// wenn Login-Button gedrückt und POST ausgeführt wurde
		if r.PostForm.Has("login") && r.Method == http.MethodPost {
			// Eingabefelder (username und password) auslesen
			password := []byte(r.PostFormValue("passwd"))
			username := r.PostFormValue("uname")
			// user authentifizieren
			successful := AuthenticateUser(username, password)
			// user erfolgreich authentifiziert
			if successful {
				// neue session erstellen
				sessionToken, expires := createSession(username)
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// übermitteltes Formular parsen
	r.ParseForm()
	// wenn Register-Button gedrückt und POST ausgeführt wurde
	if r.PostForm.Has("register") && r.Method == http.MethodPost {
		// exisitiert der Nutzername schon?
		duplicate := isDuplicateUsername(r.PostFormValue("uname"))
		// wenn Nutzername schon exisitert
		if duplicate {
			// Response header 401 setzen
			w.WriteHeader(http.StatusUnauthorized)
			// Fehlermeldung für Nutzer anzeigen
			templates.TempError.Execute(w, error2.CreateError(error2.DuplicateUserName, "/register"))
			return
			// Nutzername exisitert noch nicht --> register möglich
		} else {
			// Passwort hashen
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("passwd")), bcrypt.DefaultCost)
			// username auslesen
			username := r.PostFormValue("uname")
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
			sessionToken, expires := createSession(username)
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
			templates.TempError.Execute(w, error2.CreateError(error2.Authentification, "/"))
			return
		}
	}
}

func refreshCookie(r *http.Request) (sessionToken string, expires time.Time) {
	// Cookie auslesen
	cookie, _ := r.Cookie("session_token")
	// Sessiontoken auslesen
	sessionToken = cookie.Value
	// session auslesen
	session, _ := sessions[sessionToken]
	// Session ist valide, da zuvor CheckCookie ausgeführt wurde
	// expires um 10 min verlägern
	session.expires = session.expires.Add(1 * time.Minute)
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

func checkCookie(r *http.Request) (successful bool) {
	// Cookie auslesen
	cookie, err := r.Cookie("session_token")
	// kein Cookie
	if err == http.ErrNoCookie {
		return false
	}
	// Sessiontoken auslesen
	sessionToken := cookie.Value
	// session auslesen
	session, ok := sessions[sessionToken]
	// keine Session zu Sessiontoken gefunden
	if !ok {
		return false
	}
	// SessionToken is abgelaufen
	if session.isExpired() {
		// Session löschen
		delete(sessions, sessionToken)
		return false
	}
	return true
}

func isDuplicateUsername(username string) (isDuplicate bool) {
	// existiert der username schon?
	_, ok := users[username]
	return ok
}

func createSession(username string) (sessionToken string, expires time.Time) {
	// Sessiontoken generieren
	sessionToken = createUUID(25)
	// Session läuft nach x Minuten ab
	expires = time.Now().Add(1 * time.Minute)
	// Session anhand des Sessiontokens speichern
	sessions[sessionToken] = &session{
		uname:   username,
		expires: expires,
	}
	return sessionToken, expires
}
