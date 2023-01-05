package data

import (
	"time"
)

// TimeSeries - Substruct of Appointment
type TimeSeries struct {
	Repeat    bool `json:"repeat"`
	Intervall int  `json:"intervall"`
}

// Share - Substruct of Appointment
type Share struct {
	Public bool `json:"public"`
	Tokens []string
	Voting []bool
}

// Appointment - Information of an Appointment
type Appointment struct {
	Id            int        `json:"id"`
	DateTimeStart time.Time  `json:"dateTimeStart"`
	DateTimeEnd   time.Time  `json:"dateTimeEnd"`
	Location      string     `json:"location"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Userid        int        `json:"userid"`
	Timeseries    TimeSeries `json:"timeseries"`
	Share         Share      `json:"share"`
}

// User - Information of an User
type User struct {
	UserName           string `json:"userName"`
	Password           string `json:"password"`
	UserLevel          int    `json:"userLevel"`
	Id                 int    `json:"id"`
	Appointments       map[int]Appointment
	SharedAppointments map[string][]Appointment
}

// NewUser constructs a new User
func NewUser(name, pw string, id, userLevel int) User {
	return User{name, pw, userLevel, id, make(map[int]Appointment), make(map[string][]Appointment)}
}

// NewAppointment constructs a new Appointment
func NewAppointment(title, description, location string, dateTimeStart, dateTimeEnd time.Time, id, userId int, repeat bool, intervall int, public bool) Appointment {
	res := Appointment{id, dateTimeStart, dateTimeEnd, location, title, description, userId, TimeSeries{repeat, intervall}, Share{public, make([]string, 0), make([]bool, 0)}}
	return res
}

// GetDescriptionFromInterval returns a string representing the information of an Appointment
func (ap Appointment) GetDescriptionFromInterval() string {
	switch ap.Timeseries.Intervall {
	case 1:
		return "täglich"
	case 7:
		return "wöchentlich"
	case 30:
		return "monatlich"
	case 365:
		return "jährlich"
	default:
		return "keine"
	}
}
