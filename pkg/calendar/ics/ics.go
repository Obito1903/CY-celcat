package ics

import (
	"log"
	"os"
	"time"

	"github.com/Obito1903/CY-celcat/pkg/calendar"
	iCal "github.com/arran4/golang-ical"
)

func addEventToICS(icsCal *iCal.Calendar, event calendar.Event) {
	icsEvent := icsCal.AddEvent(event.Id)
	icsEvent.SetCreatedTime(time.Now())
	icsEvent.SetDtStampTime(time.Now())
	icsEvent.SetModifiedAt(time.Now())
	icsEvent.SetStartAt(event.Start)
	icsEvent.SetEndAt(event.End)

	summury := event.Category
	if len(event.Subjects) > 0 {
		summury += " - "
		for idx, subjects := range event.Subjects {
			if idx > 1 {
				summury += ", "
			}
			summury += subjects
		}
	}
	icsEvent.SetSummary(summury)
	if len(event.Location) > 0 {
		var locations string
		for idx, location := range event.Location {
			if idx > 0 {
				locations += ", "
			}
			locations += location
		}
		icsEvent.SetLocation(locations)
	}

	var organizer string
	if len(event.Professors) > 0 {
		for idx, professor := range event.Professors {
			if idx > 0 {
				professor += ", "
			}
			organizer += professor
		}
		icsEvent.SetOrganizer(organizer)
	}
	icsEvent.SetDescription(organizer + " | " + summury + event.Notes)
}

func CalendarToICS(cal calendar.Calendar) *iCal.Calendar {
	icsCal := iCal.NewCalendar()
	for _, event := range cal.Events {
		addEventToICS(icsCal, event)
	}
	return icsCal
}

func IcsToFile(cal *iCal.Calendar, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("Could not save ICS.", err)
	}
	defer f.Close()
	_, err = f.WriteString(cal.Serialize())
	if err != nil {
		log.Fatal("Could not save ICS.", err)
	}
}
