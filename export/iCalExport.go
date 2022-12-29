package export

import (
	"go_cal/data"
	"time"
)

type RRule struct {
	Freq    string
	ByDay   string
	ByMonth int
}

type VTimeZoneStandard struct {
	DTStart      time.Time
	RRule        RRule //Maybe split to own struct
	TZOffSetFrom float64
	TZOffSetTo   float64
}
type VTimeZoneDaylight struct {
	DTStart      time.Time
	RRule        RRule
	TZOffSetFrom float64
	TZOffSetTo   float64
}

type VTimeZone struct {
	TZID     string
	Standard VTimeZoneStandard
	Daylight VTimeZoneDaylight
}

type VEvent struct {
	UID         string
	Organizer   string
	Location    string
	Geo         string
	Summary     string
	Description string
	Class       string
	DTStart     time.Time
	DTEnd       time.Time
	DTStamp     time.Time
}

type ICal struct {
	Version float64
	VTimeZone
	VEvent
}

func NewRRule() RRule {

}

func NewVTTimeZoneStandard() VTimeZoneStandard {

}

func NewVTimeZoneDaylight() VTimeZoneDaylight {

}

func NewVTimeZone() VTimeZone {

}

func NewVEvent(ap *data.Appointment) VEvent {

}

func NewICal(ap *data.Appointment) ICal {

	//return ICal{2.0}
}

func (ics *ICal) ToString() string {

}
