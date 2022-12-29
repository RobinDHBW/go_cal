package dataModel

import (
	"encoding/json"
	"go_cal/data"
	"go_cal/fileHandler"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func encryptPW(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatalln(err)
	}
	return string(hash)
}

type DataModel struct {
	UserMap map[int]data.User
	fH      fileHandler.FileHandler
}

func NewDM(dataPath string) DataModel {
	fH := fileHandler.NewFH(dataPath)
	sList := fH.ReadAll()
	uMap := make(map[int]data.User)

	for _, uString := range sList {
		var user data.User
		json.Unmarshal([]byte(uString), &user)
		uMap[user.Id] = user
	}

	return DataModel{uMap, fH}
}

func (dm *DataModel) GetUserById(id int) *data.User {
	if res, ok := dm.UserMap[id]; ok {
		return &res
	}
	return nil
}

func (dm *DataModel) AddUser(name, pw string, userLevel int) *data.User {
	user := data.NewUser(name, encryptPW(pw), len(dm.UserMap)+1, userLevel)

	write, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	//@TODO make parallel
	dm.UserMap[user.Id] = user
	dm.fH.SyncToFile(write, user.Id)
	return &user
}

// Call by reference or call by value?
func (dm *DataModel) AddAppointment(id int, ap data.Appointment) *data.User {
	user := dm.GetUserById(id)
	//user.Appointments = append(user.Appointments, ap)
	user.Appointments[ap.Id] = ap
	return user
}

func (dm *DataModel) DeleteAppointment(apId, uId int) *data.User {
	user := dm.GetUserById(uId)
	delete(user.Appointments, apId)
	//res := data.NewUser("abc", "abc", 1, 1)

	return user
}

func (dm *DataModel) ComparePW(clear, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clear))
	if err != nil {
		return false
	}

	return true
}
