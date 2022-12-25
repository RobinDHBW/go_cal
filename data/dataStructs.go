package data

import "time"

type TimeSeries struct {
	Repeat    bool `json:"repeat"`
	Intervall int  `json:"intervall"`
}

type Share struct {
	Public bool   `json:"public"`
	Url    string `json:"url"`
}

type Appointment struct {
	DateTime    time.Time  `json:"dateTime"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Userid      int        `json:"userid"`
	Timeseries  TimeSeries `json:"timeseries"`
	Share       Share      `json:"share"`
}

type User struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	UserLevel    int    `json:"userLevel"`
	Id           int    `json:"id"`
	Appointments []Appointment
}

func NewUser(name, pw string, id, userLevel int) User {
	return User{name, pw, userLevel, id, nil}
}

func NewAppointment(title, description string, dateTime time.Time, userId int, repeat bool, intervall int, public bool, url string) Appointment {
	return Appointment{dateTime, title, description, userId, TimeSeries{repeat, intervall}, Share{public, url}}
}
