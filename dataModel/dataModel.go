package dataModel

import (
	"encoding/json"
	"go_cal/data"
	"go_cal/fileHandler"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//Class to hold all data and coordinate sync to/from file

func encryptPW(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatalln(err)
	}
	return string(hash)
}

type DataModel struct {
	UserList []data.User
	fH       fileHandler.FileHandler
}

func NewDM(dataPath string) DataModel {
	fH := fileHandler.NewFH(dataPath)
	sList := fH.ReadAll()
	var uList []data.User

	for _, uString := range sList {
		var user data.User
		json.Unmarshal([]byte(uString), &user)
		uList = append(uList, user)
	}

	return DataModel{uList, fH}
}

func (dm *DataModel) GetUserById(id int) *data.User {
	var res data.User
	for _, user := range dm.UserList {
		if user.Id == id {
			res = user
		}
	}
	return &res
}

func (dm *DataModel) AddUser(name, pw string, userLevel int, appointment []data.Appointment) *data.User {
	user := data.NewUser(name, encryptPW(pw), len(dm.UserList)+1, userLevel)
	if appointment != nil {
		for _, ap := range appointment {
			user = *dm.AddAppointment(user.Id, ap)
		}
	}
	write, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	//@TODO make parallel
	dm.UserList = append(dm.UserList, user)
	dm.fH.SyncToFile(write, user.Id)
	return &user
}

// Call by reference or call by value?
func (dm *DataModel) AddAppointment(id int, ap data.Appointment) *data.User {
	user := dm.GetUserById(id)
	user.Appointments = append(user.Appointments, ap)
	return user
}

//func (dm *DataModel) DeleteAppointment(id int) data.User {
//	return data.NewUser("abc", "abc", 1, 1)
//}

func (dm *DataModel) ComparePW(clear, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clear))
	if err != nil {
		return false
	}

	return true
}
