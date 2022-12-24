package data

type Appointment struct {
	DateTime    string `json:"dateTime"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Userid      int    `json:"userid"`
	UserLevel   int    `json:"userLevel"`
	Timeseries  struct {
		Repeat    bool `json:"repeat"`
		Intervall int  `json:"intervall"`
	} `json:"timeseries"`
	Share struct {
		Public bool   `json:"public"`
		Url    string `json:"url"`
	} `json:"share"`
}

type User struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	Id           int    `json:"id"`
	Appointments []Appointment
}

func NewUser(name, pw string, id int) User {
	return User{name, pw, id, nil}
}
