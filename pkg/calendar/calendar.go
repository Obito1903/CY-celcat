package calendar

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/Obito1903/CY-celcat/pkg/celcat"
)

type Event struct {
	Id         string
	Start      time.Time
	End        time.Time
	AllDay     bool
	Category   string
	Subjects   []string
	Location   []string
	Professors []string
	Notes      string
}

type Calendar struct {
	Name   string
	Events []Event
}

type ByDate []Event

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Start.Before(a[j].Start) }

func (c *Calendar) TomorrowFirstEvent() Event {
	sort.Sort(ByDate(c.Events))
	tomorrow := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1)
	for _, event := range c.Events {
		if event.Start.After(tomorrow) {
			return event
		}
	}
	return Event{}
}

func (calendar Calendar) NextEventToJson() string {
	event := calendar.TomorrowFirstEvent()
	if event.Id == "" {
		return "No events tomorrow"
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return "Error while marshalling event"
	}
	return string(jsonEvent)
}

func FirstDayOfISOWeek(date time.Time) time.Time {
	for date.Weekday() != time.Monday {
		if (date.Weekday() == time.Sunday) || (date.Weekday() == time.Saturday) {
			date = date.AddDate(0, 0, 1)
		} else {
			date = date.AddDate(0, 0, -1)
		}
	}
	return date
}

func eventFromCelcat(celcatEvent celcat.CelcatCalEvent) Event {
	var event Event
	event.Id = celcatEvent.Id
	event.Start, _ = time.ParseInLocation("2006-01-02T15:04:05", celcatEvent.Start, time.Local)
	event.End, _ = time.ParseInLocation("2006-01-02T15:04:05", celcatEvent.End, time.Local)
	event.AllDay = celcatEvent.AllDay
	for _, element := range celcatEvent.Elements {
		switch element.Label {
		case "Catégorie":
			event.Category = element.Content
		case "Matière", "Name":
			event.Subjects = append(event.Subjects, element.Content)
		case "Notes":
			event.Notes = element.Content
		default:
			switch element.EntityType {
			case 102:
				event.Location = append(event.Location, element.Content)
			case 101:
				event.Professors = append(event.Professors, element.Content)
			}
		}
	}
	return event
}

func FromCelcat(celcatCalendar []celcat.CelcatCalEvent, name string) Calendar {
	var calendar Calendar
	calendar.Name = name
	for _, event := range celcatCalendar {
		calEvent := eventFromCelcat(event)
		if event.AllDay {
			continue
		}
		calendar.Events = append(calendar.Events, calEvent)
	}
	return calendar
}
