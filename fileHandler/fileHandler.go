package fileHandler

import (
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

type Appointment struct {
	DateTime    string `json:"dateTime"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Userid      int    `json:"userid"`
	UserLevel   int    `json:"userLevel"`
	Timeseries  struct {
		Repeat    bool `json:"repeat"`
		Intervall int  `json:"intervall"`
	} `json:"timeseries"`
	Share struct {
		Public bool   `json:"public"`
		Url    string `json:"url"`
	} `json:"share"`
}

type User struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	Id           int    `json:"id"`
	Appointments []Appointment
}

func NewUser(name, pw string, id int) User {
	return User{name, pw, id, nil}
}

type FileHandler struct {
	dataPath  string
	fileNames []string
}

// Initialize structs from disk
func NewFH(dataPath string) FileHandler {
	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	fN := make([]string, len(files))
	for _, f := range files {
		fN = append(fN, f.Name())
	}
	return FileHandler{dataPath, fN}
}

func (fh FileHandler) SyncToFile(json []byte, id int) {
	fN := strconv.Itoa(id) + ".json"
	file, err := os.Create(path.Join(fh.dataPath, fN))
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	file.Write(json)
}

func (fh FileHandler) ReadFromFile(id int) string {
	fP := path.Join(fh.dataPath, strconv.Itoa(id)+".json")
	file, err := os.Open(fP)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//var user User

	byteVal, _ := io.ReadAll(file)
	return string(byteVal)
	//json.Unmarshal(byteVal, &user)
	//return user
}
