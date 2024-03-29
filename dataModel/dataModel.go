// Matrikelnummern:
// 9495107, 4706893, 9608900

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

// InitDataModel
// creates a new DataModel and declares the var Dm
func InitDataModel(path string) {
	Dm = NewDM(path)
}

// encryptPW
// returns the salted-hash encrypted string of given clear passwordstring
func encryptPW(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

// DataSync
// syncs a user
func dataSync(user *data.User, dm *DataModel) {
	write, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	dm.UserMap[user.Id] = *user
	dm.fH.SyncToFile(write, user.Id)
}

// DataModel
// represents an Object to store and handle all Information
type DataModel struct {
	UserMap map[int]data.User
	fH      fileHandler.FileHandler
}

// NewDM
// constructs a new DataModel instance
func NewDM(dataPath string) DataModel {
	apID = 0
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

// GetUserById
// returns a Pointer to a User by given id
func (dm *DataModel) GetUserById(id int) *data.User {
	if res, ok := dm.UserMap[id]; ok {
		return &res
	}
	return nil
}

// GetUserByName
// return a Pointer to a User by given Name
func (dm *DataModel) GetUserByName(search string) *data.User {
	for _, val := range dm.UserMap {
		if val.UserName == search {
			return &val
		}
	}
	return nil
}

// AddUser
// adds a new User to DataModel
// returns a Pointer to it
func (dm *DataModel) AddUser(name, pw string, userLevel int) (*data.User, error) {
	for _, val := range dm.UserMap {
		if val.UserName == name {
			return nil, errors.New("username already exists")
		}
	}
	user := data.NewUser(name, encryptPW(pw), len(dm.UserMap)+1, userLevel)

	dataSync(&user, dm)
	return &user, nil
}

// AddAppointment
// adds a new Appointment to a user
// returns a Pointer to the User and the Appointment
func (dm *DataModel) AddAppointment(userId int, title, description, location string, dateTimeStart, dateTimeEnd time.Time, repeat bool, intervall int, public bool) (*data.User, *data.Appointment) {
	apID++
	ap := data.NewAppointment(title, description, location, dateTimeStart, dateTimeEnd, apID, userId, repeat, intervall, public)

	user := dm.GetUserById(userId)
	user.Appointments[ap.Id] = ap

	dataSync(user, dm)
	return user, &ap
}

// AddSharedAppointment
// add SharedAppointment for user specified by userId
// if title already exists, appointment is added as additional date for SharedAppointment
// Tokens of new appointment are equal to tokens from existing, if one exists
// Voting results are initially set to false, return user
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

	dataSync(user, dm)
	return user
}

// DeleteAppointment deletes an Appointment by given id
// returns a Pointer to the corresponding User
func (dm *DataModel) DeleteAppointment(apId, uId int) *data.User {
	user := dm.GetUserById(uId)
	if user == nil {
		log.Fatal("No user for id: ", uId)
		return nil
	}
	delete(user.Appointments, apId)

	dataSync(user, dm)
	return user
}

// EditAppointment overwrites the given Appointment
// returns a Pointer to the corresponding User
func (dm *DataModel) EditAppointment(uId int, ap *data.Appointment) *data.User {
	user := dm.GetUserById(uId)
	user.Appointments[ap.Id] = *ap

	dataSync(user, dm)
	return user
}

// GetAppointmentsBySearchString searches in title and description
// returns a Pointer to the corresponding User and a map of the matching Appointments
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

// ComparePW compares a plaintext and a hash
// returns as a matching result a boolean value
func (dm *DataModel) ComparePW(clear, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clear))
	if err != nil {
		return false
	}

	return true
}

// GetAppointmentsForUser
// returns a Map of all Appointments of a User
func (dm *DataModel) GetAppointmentsForUser(uId int) *map[int]data.Appointment {
	user := dm.GetUserById(uId)
	return &user.Appointments
}

// AddTokenToSharedAppointment
// Add token for invited user "username" for SharedAppointment "title" for user with id "id"
// adds token for all dates of shared appointment, sets votes initially to false
// return error if invited user is already invited
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
	dataSync(user, dm)
	return nil
}

// SetVotingForToken
// sets voting result of termin proposals to corresponding shared Termin
// return error if voting is not allowed
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
		for i := range user.SharedAppointments[title] {
			user.SharedAppointments[title][i].Share.Voting[index] = false
		}
		// Votes bei zugesagten Terminen auf true setzen
		for i := range votes {
			user.SharedAppointments[title][votes[i]].Share.Voting[index] = true
		}
		dataSync(user, dm)
		return nil
	} else {
		return errors.New("voting not allowed")
	}
}

// IsVotingAllowed
// checks whether URL with provided Query Parameters is allowed for voting
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

// CreateURL
// creates URL with 4 query parameters: username, title, token and invitor
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

// createToken
// generates a random token with length n
func createToken(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// InitSeed
// initialize rand
func InitSeed() {
	rand.Seed(time.Now().UnixNano())
}

// DeleteSharedAppointment
// deletes shared appointment with title "title" from user with "uId"
// return user
func (dm *DataModel) DeleteSharedAppointment(title string, uId int) *data.User {
	user := dm.GetUserById(uId)
	delete(user.SharedAppointments, title)

	dataSync(user, dm)
	return user
}
