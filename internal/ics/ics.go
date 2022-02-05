package ics

import (
	"fmt"
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

	if len(event.Location) > 0 {
		var organizer string
		for idx, professor := range event.Professors {
			if idx > 0 {
				professor += ", "
			}
			organizer += professor
		}
		icsEvent.SetOrganizer(organizer)
	}
}

func CalendarToICS(cal calendar.Calendar) {
	icsCal := iCal.NewCalendar()
	for _, event := range cal.Events {
		addEventToICS(icsCal, event)
	}
	fmt.Println(icsCal.Serialize())
}
