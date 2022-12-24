package authentication

import (
	"encoding/json"
	"fmt"
	"go_cal/calendarView"
	"go_cal/templates"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var users = map[string][]byte{}

type Credentials struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type session struct {
	uname   string
	expires time.Time
}

var sessions = map[string]session{}

func (s session) isExpired() bool {
	return s.expires.Before(time.Now())
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// zu überprüfen: Datei zu User finden (exisitert der User?), Passwort abgleichen, Cookie setzen

	r.ParseForm()
	if r.PostForm.Has("login") && r.Method == http.MethodPost {
		password := []byte(r.PostFormValue("passwd"))
		username := r.PostFormValue("uname")
		successful := AuthenticateUser(username, password)
		if successful {
			sessionToken := createUUID(10)
			expires := time.Now().Add(120 * time.Second)
			sessions[sessionToken] = session{
				uname:   username,
				expires: expires,
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expires,
			})

			http.Redirect(w, r, "/updateCalendar", http.StatusFound)
			return
		} else {
			r.Method = http.MethodGet
			http.Redirect(w, r, "error?type=authentification&link="+url.QueryEscape("/"), http.StatusContinue)
			return
		}
	}
	templates.TempLogin.Execute(w, nil)

	//if r.PostFormValue("uname") == "" && r.PostFormValue("passwd") == "" {
	//	templates.TempLogin.Execute(w, nil)
	//} else {
	//	file, err := os.ReadFile("./files/" + r.PostFormValue("uname") + ".json")
	//	if err != nil {
	//		fmt.Println(err)
	//		templates.TempLogin.Execute(w, nil)
	//		return
	//	}
	//	var user Credentials
	//	json.Unmarshal(file, &user)
	//	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(r.PostFormValue("passwd")))
	//	if err == nil {
	//		fmt.Println("Successfully logged in")
	//		templates.TempInit.Execute(w, calendarView.Cal)
	//	} else {
	//		fmt.Println("login failed")
	//		templates.TempLogin.Execute(w, nil)
	//	}
	//}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// exisitiert der Nutzer schon?

	if r.PostFormValue("uname") == "" && r.PostFormValue("passwd") == "" {
		templates.TempRegister.Execute(w, nil)
	} else {
		passwd, _ := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("passwd")), bcrypt.DefaultCost)
		user := Credentials{
			Username: r.PostForm.Get("uname"),
			Password: passwd,
		}

		// write userinfo to filesystem

		// cookie setzen

		folder, err := os.Open("./files")
		if err != nil {
			fmt.Println(err)
			return
		}

		files, err := folder.Readdir(0)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range files {
			if user.Username == strings.Split(file.Name(), ".")[0] {
				templates.TempRegister.Execute(w, nil)
				return
			}
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

		templates.TempInit.Execute(w, calendarView.Cal)
	}
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

func AuthenticateUser(username string, password []byte) (successful bool) {
	expectedPasswordHash, ok := users[username]

	if !ok || bcrypt.CompareHashAndPassword(expectedPasswordHash, password) != nil {
		return false
	} else {
		return true
	}
}

//func createUser(uname *string, passwd *string) error {
//
//}
//
//func createCookie() {
//
//}
