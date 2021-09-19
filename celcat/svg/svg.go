package svg

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/Obito1903/CY-celcat/celcat/common"
	svg "github.com/ajstarks/svgo"
)

func makeGrid(canvas *svg.SVG, width int, height int, day_start int, day_end int) int {
	day_length := day_end - day_start // End of the Day minus Start of the day
	nb_line := day_length * 2 // Number of line to draw
	leftPanel := width / 25
	eventListwidth := width - leftPanel
	eventWidth := eventListwidth / 5
	for i := 0; i < 5; i++ {
		canvas.Line(leftPanel+eventWidth*i, height, leftPanel+eventWidth*i, 0, "stroke:white")
	}
	timeLineHeight := (height / nb_line)
	currentTimeLine := time.Date(2021, 1, 1, day_start, 0, 0, 0, time.UTC)
	gapBetweenLine, _ := time.ParseDuration("30m")
	for i := 0; i < nb_line+1; i++ {
		canvas.Text(0, (timeLineHeight)*i+20, fmt.Sprintf("%dh%02d", currentTimeLine.Hour(), currentTimeLine.Minute()), "font-size:21px;font-family:Roboto;fill:white;font-size:20px")
		canvas.Line(0, (timeLineHeight)*i, width, (timeLineHeight)*i, "stroke:white")
		currentTimeLine = currentTimeLine.Add(gapBetweenLine)
	}
	return eventWidth
}

func colorEvent(eventType string) string {
	var color string

	switch eventType {
	case "TD":
		color = "#10A4C4"
	case "TP":
		color = "teal"
	case "CM":
		color = "#F04747"
	case "Examens":
		color = "purple"
	case "Tiers temps":
		color = "pink"
	default:
		color = "black"
	}
	return color
}

func addEvent(canvas *svg.SVG, width int, height int, eventWidth int, day_start int, day_end int, event common.CalEvent) {

	day_length := day_end - day_start
	minutes_height := float64(height) / float64((day_length*60))
	//Draw event box
	StartTimeOfDay := event.Start.Hour()*60 + event.Start.Minute() - day_start * 60
	StartY := int((minutes_height * (float64(StartTimeOfDay))))
	EndTimeOfDay := event.End.Hour()*60 + event.End.Minute() - day_start * 60
	eventHeight := int((minutes_height * float64((EndTimeOfDay - StartTimeOfDay))))
	StartX := (width / 25) + (int(event.Start.Weekday())-1)*eventWidth + 5
	canvas.Roundrect(StartX, StartY, eventWidth-10, eventHeight, 5, 5, "stroke:"+colorEvent(event.Category)+";stroke-width:3;fill:#32353B")

	canvas.Text(StartX+8, StartY+24, event.Module, "font-size:21px;font-family:Roboto;fill:white")
	canvas.Text(StartX+8, StartY+48, event.Prof, "font-size:21px;font-family:Roboto;fill:#40AD7B")

	startTimeStr := fmt.Sprintf("%02dh%02d", event.Start.Hour(), event.Start.Minute())
	canvas.Text(StartX+eventWidth-len(startTimeStr)*13-18, StartY+22, startTimeStr, "font-size:20px;font-family:Roboto;fill:white")

	endTimeStr := fmt.Sprintf("%02dh%02d", event.End.Hour(), event.End.Minute())
	canvas.Text(StartX+eventWidth-len(startTimeStr)*13-18, StartY+eventHeight-8, endTimeStr, "font-size:20px;font-family:Roboto;fill:white")

}

func CreateSVG(calendar common.Calendar) *bytes.Buffer {
	svgString := ""
	s := bytes.NewBufferString(svgString)
	fileWriter := bufio.NewWriter(s)
	width := 1920
	height := 1080
	canvas := svg.New(fileWriter)
	canvas.Start(width, height)
	canvas.Def()
	fmt.Fprint(canvas.Writer, "<style>\n@import url(\"https://fonts.googleapis.com/css2?family=Roboto:wght@100;300;400\");\n</style>\n")
	canvas.DefEnd()
	canvas.Rect(0, 0, width, height, "fill:#36393F")
	eventWidth := makeGrid(canvas, width, height, 8, 19)
	for _, event := range calendar {
		addEvent(canvas, width, height, eventWidth, 8, 19, event)
	}
	canvas.End()

	fileWriter.Flush()
	return s
}
