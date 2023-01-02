package export

import (
	"fmt"
	"go_cal/data"
	"time"
)

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
	VEvent
}

func NewVEvent(uid, location, summary, description, class string, dtstart, dtend, dtstamp time.Time) VEvent {
	return VEvent{uid, location, summary, description, class, dtstart, dtend, dtstamp}
}

func NewICal(ap *data.Appointment) ICal {
	vevent := NewVEvent(fmt.Sprintf("%d", ap.Id), ap.Location, ap.Title, ap.Description, "PUBLIC", ap.DateTimeStart, ap.DateTimeEnd, time.Now())
	return ICal{2.0, "Cal_App/go_cal", "PUBLISH", vevent}
}

func (ics *ICal) ToString() string {
	res := "BEGIN:VCALENDAR"
	res += "\nVERSION:" + fmt.Sprintf("%f", ics.Version)
	res += "\nPRODID:" + ics.ProdID
	res += "\nMETHOD:" + ics.Method
	res += "\nBEGIN:VEVENT"
	res += "\nUID:" + ics.VEvent.UID
	res += "\nLOCATION:" + ics.VEvent.Location
	res += "\nSUMMARY" + ics.VEvent.Summary
	res += "\nDESCRIPTION:" + ics.VEvent.Description
	res += "\nCLASS:" + ics.VEvent.Class
	res += "\nDTSTART:" + ics.VEvent.DTStart.String()
	res += "\nDTEND:" + ics.VEvent.DTEnd.String()
	res += "\nDTSTAMP" + ics.VEvent.DTStamp.String()
	res += "\nEND:VEVENT"
	res += "\nEND:VCALENDAR"

	return res
}
