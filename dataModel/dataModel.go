package dataModel

import (
	"encoding/json"
	"errors"
	"go_cal/data"
	"go_cal/fileHandler"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

var Dm DataModel
var apID = 0

func InitDataModel(path string) {
	Dm = NewDM(path)
}

func encryptPW(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func DataSync(user *data.User, dm *DataModel) {
	//@TODO make parallel
	write, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	dm.UserMap[user.Id] = *user
	dm.fH.SyncToFile(write, user.Id)
}

//func CheckDate(toCheck, from, to time.Time) bool {
//	return from.Equal(toCheck) || !from.After(toCheck) || to.Equal(toCheck) || !to.Before(toCheck)
//}

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
		err := json.Unmarshal([]byte(uString), &user)

		if err != nil {
			log.Fatal(err)
		}

		uMap[user.Id] = user
		for _, ap := range user.Appointments {
			if ap.Id > apID {
				apID = ap.Id
			}
		}
	}

	return DataModel{uMap, fH}
}

func (dm *DataModel) GetUserById(id int) *data.User {
	if res, ok := dm.UserMap[id]; ok {
		return &res
	}
	return nil
}

func (dm *DataModel) GetUserByName(search string) *data.User {
	for _, val := range dm.UserMap {
		if val.UserName == search {
			return &val
		}
	}
	return nil
}

func (dm *DataModel) AddUser(name, pw string, userLevel int) (*data.User, error) {
	for _, val := range dm.UserMap {
		if val.UserName == name {
			return nil, errors.New("username already exists")
		}
	}
	user := data.NewUser(name, encryptPW(pw), len(dm.UserMap)+1, userLevel)

	DataSync(&user, dm)
	return &user, nil
}

func (dm *DataModel) AddAppointment(userID int, title, description, location string, dateTimeStart, dateTimeEnd time.Time, userId int, repeat bool, intervall int, public bool) (*data.User, *data.Appointment) {
	apID++
	ap := data.NewAppointment(title, description, location, dateTimeStart, dateTimeEnd, apID, userId, repeat, intervall, public)

	user := dm.GetUserById(userID)
	user.Appointments[ap.Id] = ap

	DataSync(user, dm)
	return user, &ap
}

func (dm *DataModel) AddSharedAppointment(id int, title string, ap data.Appointment) *data.User {
	user := dm.GetUserById(id)
	user.SharedAppointments[title] = append(user.SharedAppointments[title], ap)

	DataSync(user, dm)
	return user
}

func (dm *DataModel) DeleteAppointment(apId, uId int) *data.User {
	user := dm.GetUserById(uId)
	delete(user.Appointments, apId)

	DataSync(user, dm)
	return user
}

func (dm *DataModel) EditAppointment(uId int, ap *data.Appointment) *data.User {
	user := dm.GetUserById(uId)
	user.Appointments[ap.Id] = *ap

	DataSync(user, dm)
	return user
}

//func (dm *DataModel) GetAppointmentsByTimeFrame(uId int, tFrom, tTo time.Time) (*data.User, *map[int]data.Appointment) {
//	user := dm.GetUserById(uId)
//	res := make(map[int]data.Appointment)
//	for key, val := range user.Appointments {
//		if CheckDate(val.DateTimeStart, tFrom, tTo) || CheckDate(val.DateTimeEnd, tFrom, tTo) {
//			res[key] = val
//		}
//	}
//
//	return user, &res
//}

func (dm *DataModel) GetAppointmentsBySearchString(uId int, search string) (*data.User, *map[int]data.Appointment) {
	user := dm.GetUserById(uId)
	res := make(map[int]data.Appointment)
	for key, val := range user.Appointments {
		if strings.Contains(val.Title, search) || strings.Contains(val.Description, search) {
			res[key] = val
		}
	}

	return user, &res
}

func (dm *DataModel) ComparePW(clear, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clear))
	if err != nil {
		return false
	}

	return true
}

func (dm *DataModel) GetAppointmentsForUser(uId int) map[int]data.Appointment {
	user := dm.GetUserById(uId)
	return user.Appointments
}
func (dm *DataModel) AddTokenToSharedAppointment(id int, title, url string) {
	user := dm.GetUserById(id)
	for i := range user.SharedAppointments[title] {
		user.SharedAppointments[title][i].Share.Tokens = append(user.SharedAppointments[title][i].Share.Tokens, url)
		user.SharedAppointments[title][i].Share.Voting = append(user.SharedAppointments[title][i].Share.Voting, false)
	}

	DataSync(user, dm)
}
