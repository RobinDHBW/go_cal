package dataModel

import (
	"encoding/json"
	"errors"
	"go_cal/data"
	"go_cal/fileHandler"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	url2 "net/url"
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

func (dm *DataModel) AddAppointment(userId int, title, description, location string, dateTimeStart, dateTimeEnd time.Time, repeat bool, intervall int, public bool) (*data.User, *data.Appointment) {
	apID++
	ap := data.NewAppointment(title, description, location, dateTimeStart, dateTimeEnd, apID, userId, repeat, intervall, public)

	user := dm.GetUserById(userId)
	user.Appointments[ap.Id] = ap

	DataSync(user, dm)
	return user, &ap
}

func (dm *DataModel) AddSharedAppointment(userId int, title, location string, dateTimeStart, dateTimeEnd time.Time, repeat bool, intervall int, public bool) *data.User {
	apID++
	ap := data.NewAppointment(title, "", location, dateTimeStart, dateTimeEnd, apID, userId, repeat, intervall, public)
	user := dm.GetUserById(userId)
	user.SharedAppointments[title] = append(user.SharedAppointments[title], ap)

	length := len(user.SharedAppointments[title])
	if length > 1 {
		user.SharedAppointments[title][length-1].Share.Tokens = user.SharedAppointments[title][0].Share.Tokens
		user.SharedAppointments[title][length-1].Share.Public = true
		user.SharedAppointments[title][length-1].Share.Voting = make([]bool, len(user.SharedAppointments[title][0].Share.Voting))
		for i := range user.SharedAppointments[title][length-1].Share.Voting {
			user.SharedAppointments[title][length-1].Share.Voting[i] = false
		}
	}

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

func (dm *DataModel) AddTokenToSharedAppointment(id int, title, url, username string) error {
	user := dm.GetUserById(id)
	for _, val := range user.SharedAppointments[title] {
		for _, text := range val.Share.Tokens {
			parsedUrl, _ := url2.Parse(text)
			if username == parsedUrl.Query().Get("username") {
				return errors.New("duplicate invited username")
			}
		}
	}
	for i := range user.SharedAppointments[title] {
		user.SharedAppointments[title][i].Share.Tokens = append(user.SharedAppointments[title][i].Share.Tokens, url)
		user.SharedAppointments[title][i].Share.Voting = append(user.SharedAppointments[title][i].Share.Voting, false)
	}
	DataSync(user, dm)
	return nil
}

func (dm *DataModel) SetVotingForToken(user *data.User, votes []int, title, token, username string) error {
	if IsVotingAllowed(title, token, user, username) {
		var index int
		for i, text := range user.SharedAppointments[title][0].Share.Tokens {
			if strings.Contains(text, token) && strings.Contains(text, username) {
				index = i
				break
			}
		}
		// alle Votes auf false setzen
		for i, _ := range user.SharedAppointments[title] {
			user.SharedAppointments[title][i].Share.Voting[index] = false
		}
		for i := range votes {
			user.SharedAppointments[title][votes[i]].Share.Voting[index] = true
		}
		DataSync(user, dm)
		return nil
	} else {
		return errors.New("voting not allowed")
	}
}

func IsVotingAllowed(title, token string, user *data.User, username string) bool {
	if user == nil {
		return false
	}
	query := "/terminVoting?invitor=" + user.UserName + "&termin=" + title + "&token=" + token + "&username=" + username
	if _, ok := user.SharedAppointments[title]; !ok {
		return false
	}
	for _, val := range user.SharedAppointments[title][0].Share.Tokens {
		if val == query {
			return true
		}
	}
	return false
}

func CreateURL(username, title, invitor string) string {
	token := createToken(20)
	params := url2.Values{}
	params.Add("username", username)
	params.Add("termin", title)
	params.Add("token", token)
	params.Add("invitor", invitor)
	baseUrl, _ := url2.Parse("/terminVoting")
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String()
}

func createToken(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func InitSeed() {
	rand.Seed(time.Now().UnixNano())
}
