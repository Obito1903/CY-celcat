package html

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"time"

	config "github.com/Obito1903/CY-celcat/pkg"
	"github.com/Obito1903/CY-celcat/pkg/calendar"
)

type htmlEvent struct {
	Event    calendar.Event
	Top      float32
	Height   float32
	TimeSpan string
}

type htmlCalendar struct {
	Horaires []string
	Days     map[time.Weekday][]htmlEvent
	MaxStart time.Time
	MaxEnd   time.Time
}

// Calculate the duration between the start of the day (00h00) and the current time in that day
func hourInDayToDuration(d time.Time) time.Duration {
	return time.Duration(d.Hour()*int(time.Hour) + d.Minute()*int(time.Minute))
}

// Find the earliest and latest event in the calendar then split the duration of the day into 30min chunk and add them to the Horaire array
// For a calendar starting at 7h and ending at 10h it will generate a list of timstamp like this :
// 7h | 7h30 | 8h | 8h30 | 9h | 9h30
func (htmlCal *htmlCalendar) calcHoraires(cal calendar.Calendar, week time.Time) {
	// Init min and max timestamps
	htmlCal.MaxStart = week.Add(time.Hour * 24)
	htmlCal.MaxEnd = week

	// For each event in the calendar, if it is in the provided week, compare the time it start and end with MaxStart and MaxEnd
	for _, event := range cal.Events {
		if event.End.Before(week.Add(7*24*time.Hour)) && event.Start.After(week) {
			start := time.Date(week.Year(), week.Month(), week.Day(), event.Start.Hour(), event.Start.Minute(), 0, 0, time.Local)
			end := time.Date(week.Year(), week.Month(), week.Day(), event.End.Hour(), event.End.Minute(), 0, 0, time.Local)
			if end.After(htmlCal.MaxEnd) {
				htmlCal.MaxEnd = end.Round(time.Minute * 30) // Round the timestamp to avoid having wierd calendar starting at odd minutes
			}
			if start.Before(htmlCal.MaxStart) {
				htmlCal.MaxStart = start.Round(time.Minute * 30).Add(-time.Minute * 30) // Same as above
			}
		}
	}

	// Split the duration of the days into 30min chunks
	htmlCal.MaxStart = time.Date(htmlCal.MaxStart.Year(), htmlCal.MaxStart.Month(), htmlCal.MaxStart.Day(), htmlCal.MaxStart.Hour(), (htmlCal.MaxStart.Minute()/10)*10, 0, 0, time.Local)
	for h := htmlCal.MaxStart; h.Unix() < htmlCal.MaxEnd.Unix(); h = h.Add(time.Minute * 30) {
		htmlCal.Horaires = append(htmlCal.Horaires, h.Format("15h04"))
	}
}

// Convert a standard event into a htmlEvent by adding placing and display info
func (htmlCal htmlCalendar) calEventToHtmlEvent(event calendar.Event) htmlEvent {
	dayLength := time.Duration(htmlCal.MaxEnd.Unix()-htmlCal.MaxStart.Unix()) * 1000 * 1000 * 1000
	// Start timestamp of the event relative to the start of the day in the curent week
	eventStarInDay := time.Duration(hourInDayToDuration(event.Start) - hourInDayToDuration(htmlCal.MaxStart))
	eventLength := time.Duration((event.End.Unix() - event.Start.Unix()) * 1000 * 1000 * 1000)
	htmlEv := htmlEvent{
		Event:    event,
		Top:      (float32(eventStarInDay) / float32(dayLength)) * 100,
		Height:   (float32(eventLength) / float32(dayLength)) * 100,
		TimeSpan: event.Start.Format("15h04") + "-" + event.End.Format("15h04"),
	}

	return htmlEv
}

// Convert a Calendar into an htmlCalendar
func CalToHtmlCal(cal calendar.Calendar, week time.Time) htmlCalendar {
	var htmlCal htmlCalendar
	// Init htmlCalendar according to the given calendar
	week = calendar.FirstDayOfISOWeek(week)
	htmlCal.calcHoraires(cal, week)

	// Init the days map
	htmlCal.Days = make(map[time.Weekday][]htmlEvent)
	for _, event := range cal.Events {
		if event.End.Before(week.Add(7*24*3600*1000*1000*1000)) && event.Start.After(week) {
			htmlCal.Days[event.Start.Weekday()] = append(htmlCal.Days[event.Start.Weekday()], htmlCal.calEventToHtmlEvent(event))
		}
	}
	return htmlCal
}

func (cal htmlCalendar) ToHtml(templatePath string) string {
	t, err := template.New("calendar").ParseFiles(templatePath)
	if err != nil {
		log.Fatal("Could not load template.", err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, cal)
	if err != nil {
		log.Fatal("Could not execute template.", err)
	}
	return buf.String()
}

func (htmlCal htmlCalendar) ToFile(templatePath string, outPath string) {
	html := htmlCal.ToHtml(templatePath)
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatal("Could not save HTML.", err)
	}
	defer f.Close()
	_, err = f.WriteString(html)
	if err != nil {
		log.Fatal("Could not save HTML.", err)
	}
}

func ToPng(config config.Config, htmlPath string, outPath string) {
	cmd := exec.Command(config.ChromePath, "--headless", "--no-sandbox", "--disable-gpu", "--screenshot="+outPath, fmt.Sprint("--window-size=", config.PNGWidth, ",", config.PNGHeigh), htmlPath)
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}
