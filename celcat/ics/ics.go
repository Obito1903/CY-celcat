package ics

import (
	"os"

	"github.com/Obito1903/CY-celcat/celcat/common"
	ics "github.com/arran4/golang-ical"
)

// Create the ICS serialized string from the list of json celcat event
func CreatICS(calendar common.Calendar) string {

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for _, calEvent := range calendar {
		event := cal.AddEvent(calEvent.Id)
		event.SetStartAt(calEvent.Start)
		event.SetEndAt(calEvent.End)
		event.SetLocation(calEvent.Location)
		event.SetOrganizer(calEvent.Prof)
		event.SetSummary(calEvent.Module)
	}
	return cal.Serialize()
}

// Save the ICS serialized string to a file
func SaveICS(icsCal string) {
	f, err := os.Create("data.ics")
	common.CheckErr(err)
	defer f.Close()
	_, err = f.WriteString(icsCal)
	common.CheckErr(err)
}
