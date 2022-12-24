package dataModel

import (
	"go_cal/fileHandler"
)

//Class to hold all data and coordinate sync to/from file

type DataModel struct {
	UserList []fileHandler.User
}

func (dm DataModel) getUser(id int) {

}

func (dm DataModel) addUser() {

}

//Call by referenceor call by value?
func (dm DataModel) addAppointment(id int, ap fileHandler.Appointment) {

}
