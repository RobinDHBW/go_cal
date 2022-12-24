package dataModel

import (
	"encoding/json"
	"go_cal/data"
	"go_cal/fileHandler"
)

//Class to hold all data and coordinate sync to/from file

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

func (dm DataModel) GetUserById(id int) data.User {
	var res data.User
	for _, user := range dm.UserList {
		if user.Id == id {
			res = user
		}
	}
	return res
}

func (dm DataModel) AddUser() {

}

// Call by reference or call by value?
func (dm DataModel) AddAppointment(id int, ap data.Appointment) {

}

func (dm DataModel) DeleteAppointment(id int) {

}
