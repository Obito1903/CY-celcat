package html

import (
	"fmt"
	"testing"
	"time"

	"github.com/Obito1903/CY-celcat/pkg/calendar"
)

func TestCalcHoraires(t *testing.T) {
	cal := calendar.Calendar{
		Name: "GIG1",
		Events: []calendar.Event{
			{
				Id:         "Salut",
				Start:      time.Now(),
				End:        time.Now().Add(time.Hour),
				AllDay:     false,
				Category:   "CM",
				Subjects:   []string{"Algo"},
				Location:   []string{"102"},
				Professors: []string{"Ines"},
				Notes:      "nothing",
			},
			{
				Id:         "Salut2",
				Start:      time.Now().Add(time.Hour),
				End:        time.Now().Add(time.Hour * 2),
				AllDay:     false,
				Category:   "CM",
				Subjects:   []string{"Anglais"},
				Location:   []string{"102"},
				Professors: []string{"Magda"},
				Notes:      "nothing",
			},
		},
	}
	var htmlCal htmlCalendar

	calcHoraires(cal, time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local), &htmlCal)
	fmt.Println(htmlCal.MaxEnd.Hour())
}
