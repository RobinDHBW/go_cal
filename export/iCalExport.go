package export

import (
	"fmt"
	"go_cal/data"
	"strconv"
	"time"
)

func timeToCoalString(t time.Time) string {
	return strconv.Itoa(t.Year()) + strconv.Itoa(int(t.Month())) + strconv.Itoa(t.Day()) + "T" + strconv.Itoa(t.Hour()) + strconv.Itoa(t.Minute()) + strconv.Itoa(t.Second()) + "Z"
}

type VEvent struct {
	UID         string
	Location    string
	Summary     string
	Description string
	Class       string
	DTStart     time.Time
	DTEnd       time.Time
	DTStamp     time.Time
}

type ICal struct {
	Version float64
	ProdID  string
	Method  string
	VEvent  []VEvent
}

func NewVEvent(uid, location, summary, description, class string, dtstart, dtend, dtstamp time.Time) VEvent {
	return VEvent{uid, location, summary, description, class, dtstart, dtend, dtstamp}
}

func NewICal(aps []*data.Appointment) ICal {
	vevent := make([]VEvent, 0)
	for _, ap := range aps {
		vevent = append(vevent, NewVEvent(fmt.Sprintf("%d", ap.Id), ap.Location, ap.Title, ap.Description, "PUBLIC", ap.DateTimeStart, ap.DateTimeEnd, time.Now()))
	}
	return ICal{2.0, "Cal_App/go_cal", "PUBLISH", vevent}
}

func (ics *ICal) ToString() string {
	res := "BEGIN:VCALENDAR"
	//res += "\nVERSION:" + fmt.Sprintf("%f", ics.Version)
	res += "\nVERSION:" + strconv.FormatFloat(ics.Version, 'f', 1, 64)
	res += "\nPRODID:" + ics.ProdID
	res += "\nMETHOD:" + ics.Method

	for _, event := range ics.VEvent {
		res += "\nBEGIN:VEVENT"
		res += "\nUID:" + event.UID
		res += "\nLOCATION:" + event.Location
		res += "\nSUMMARY:" + event.Summary
		res += "\nDESCRIPTION:" + event.Description
		res += "\nCLASS:" + event.Class
		res += "\nDTSTART:" + timeToCoalString(event.DTStart)
		res += "\nDTEND:" + timeToCoalString(event.DTEnd)
		res += "\nDTSTAMP:" + timeToCoalString(event.DTStamp)
		res += "\nEND:VEVENT"
	}
	res += "\nEND:VCALENDAR"

	return res
}
