package authentication

import (
	"encoding/json"
	"fmt"
	"go_cal/calendarView"
	"go_cal/templates"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

type User struct {
	Uname  string `json:"username"`
	Passwd []byte `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("uname") == "" && r.PostFormValue("passwd") == "" {
		templates.TempLogin.Execute(w, nil)
	} else {
		file, err := os.ReadFile("./files/" + r.PostFormValue("uname") + ".json")
		if err != nil {
			fmt.Println(err)
			templates.TempLogin.Execute(w, nil)
			return
		}
		var user User
		json.Unmarshal(file, &user)
		err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(r.PostFormValue("passwd")))
		if err == nil {
			fmt.Println("Successfully logged in")
			templates.TempInit.Execute(w, calendarView.Cal)
		} else {
			fmt.Println("login failed")
			templates.TempLogin.Execute(w, nil)
		}
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("uname") == "" && r.PostFormValue("passwd") == "" {
		templates.TempRegister.Execute(w, nil)
	} else {
		passwd, _ := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("passwd")), bcrypt.DefaultCost)
		user := User{
			Uname:  r.PostForm.Get("uname"),
			Passwd: passwd,
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
			if user.Uname == strings.Split(file.Name(), ".")[0] {
				templates.TempRegister.Execute(w, nil)
				return
			}
		}

		text, _ := json.Marshal(user)

		file, err := os.Create("./files/" + user.Uname + ".json")
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
