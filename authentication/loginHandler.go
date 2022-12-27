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

// map with username and corresponding hashed password
var users = map[string][]byte{}

// Credentials struct for a user
type Credentials struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

// session consist of an user and an expire time
type session struct {
	uname   string
	expires time.Time
}

// map with SessionTokens and corresponding sessions
var sessions = map[string]session{}

// determines if a Session is expired
func (s session) isExpired() bool {
	return s.expires.Before(time.Now())
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie überprüfen
	isCookieValid := CheckCookie(r)

	// no cookie sent with the request --> login-procedure
	if !isCookieValid {
		// zu überprüfen: Datei zu User finden (exisitert der User?), Passwort abgleichen, Cookie setzen
		r.ParseForm()
		if r.PostForm.Has("login") && r.Method == http.MethodPost {
			password := []byte(r.PostFormValue("passwd"))
			username := r.PostFormValue("uname")
			successful := AuthenticateUser(username, password)
			if successful {
				sessionToken, expires := createSession(username)
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: expires,
				})
				http.Redirect(w, r, "/updateCalendar", http.StatusFound)
				return
			} else {
				error2.CreateError(error2.WrongCredentials, "/", w, http.StatusUnauthorized)
				return
			}
		}
		// Cookie sent with the request + valid
	} else {
		http.Redirect(w, r, "/updateCalendar", http.StatusFound)
		return
	}
	templates.TempLogin.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm.Has("register") && r.Method == http.MethodPost {
		// exisitiert der Nutzer schon?
		duplicate := isDuplicateUsername(r.PostFormValue("uname"))
		if duplicate {
			error2.CreateError(error2.DuplicateUserName, "/register", w, http.StatusUnauthorized)
			return
		} else {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("passwd")), bcrypt.DefaultCost)
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
			// Cookie erstellen
			sessionToken, expires := createSession(username)
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expires,
			})
			// sync file to variables
			LoadUsersFromFiles()
			http.Redirect(w, r, "/updateCalendar", http.StatusFound)
			return
		}
	}
	templates.TempRegister.Execute(w, nil)
}

// TODO
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

func CheckCookie(r *http.Request) (successful bool) {
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		return false
	}

	sessionToken := cookie.Value

	// look up in session map
	session, ok := sessions[sessionToken]
	// no SessionToken found
	if !ok {
		return false
	}
	// SessionToken is expired
	if session.isExpired() {
		delete(sessions, sessionToken)
		return false
	}
	return true
}

func isDuplicateUsername(username string) (isDuplicate bool) {
	_, ok := users[username]
	return ok
}

func createSession(username string) (sessionToken string, expires time.Time) {
	sessionToken = createUUID(10)
	expires = time.Now().Add(120 * time.Second)
	sessions[sessionToken] = session{
		uname:   username,
		expires: expires,
	}
	return sessionToken, expires
}
