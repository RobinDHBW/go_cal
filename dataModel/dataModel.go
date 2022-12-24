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

func (dm DataModel) getUser(id int) {

}

func (dm DataModel) addUser() {

}

// Call by reference or call by value?
func (dm DataModel) addAppointment(id int, ap data.Appointment) {

}

func (dm DataModel) deleteAppointment(id int) {

}
