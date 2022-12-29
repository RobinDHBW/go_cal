package data

import "time"

var ApId = 0

type TimeSeries struct {
	Repeat    bool `json:"repeat"`
	Intervall int  `json:"intervall"`
}

type Share struct {
	Public bool   `json:"public"`
	Url    string `json:"url"`
}

type Appointment struct {
	Id            int        `json:"id"`
	DateTimeStart time.Time  `json:"dateTimeStart"`
	DateTimeEnd   time.Time  `json:"dateTimeEnd"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Userid        int        `json:"userid"`
	Timeseries    TimeSeries `json:"timeseries"`
	Share         Share      `json:"share"`
}

type User struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	UserLevel    int    `json:"userLevel"`
	Id           int    `json:"id"`
	Appointments map[int]Appointment
}

func NewUser(name, pw string, id, userLevel int) User {
	return User{name, pw, userLevel, id, make(map[int]Appointment)}
}

func NewAppointment(title, description string, dateTimeStart, dateTimeEnd time.Time, userId int, repeat bool, intervall int, public bool, url string) Appointment {
	res := Appointment{ApId, dateTimeStart, dateTimeEnd, title, description, userId, TimeSeries{repeat, intervall}, Share{public, url}}
	ApId++
	return res
}
