package svg

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/Obito1903/CY-celcat/celcat/common"
	svg "github.com/ajstarks/svgo"
)

func makeGrid(canvas *svg.SVG, width int, height int) int {
	leftPanel := width / 25
	eventListwidth := width - leftPanel
	eventWidth := eventListwidth / 5
	for i := 0; i < 5; i++ {
		canvas.Line(leftPanel+eventWidth*i, height, leftPanel+eventWidth*i, 0, "stroke:white")
	}
	timeLineHeight := (height / 21)
	currentTimeLine := time.Date(2021, 1, 1, 8, 0, 0, 0, time.UTC)
	gapBetweenLine, _ := time.ParseDuration("30m")
	for i := 0; i < 22; i++ {
		canvas.Text(0, (timeLineHeight)*i+60, fmt.Sprintf("%dh%02d", currentTimeLine.Hour(), currentTimeLine.Minute()), "font-family:Roboto;fill:white;font-size:20px")
		canvas.Line(0, (timeLineHeight)*i+40, width, (timeLineHeight)*i+40, "stroke:white")
		currentTimeLine = currentTimeLine.Add(gapBetweenLine)
	}
	return eventWidth
}

func colorEvent(eventType string) string {
	var color string

	switch eventType {
	case "TD":
		color = "blue"
	case "TP":
		color = "teal"
	case "CM":
		color = "red"
	case "Examens":
		color = "purple"
	case "Tiers temps":
		color = "pink"
	default:
		color = "black"
	}
	return color
}

func addEvent(canvas *svg.SVG, width int, height int, eventWidth int, event common.CalEvent) {

	//Draw event box
	StartTimeOfDay := event.Start.Hour()*60 + event.Start.Minute() - 480
	StartY := (height * (StartTimeOfDay) / 630) + 40
	EndTimeOfDay := event.End.Hour()*60 + event.End.Minute() - 480
	eventHeight := (height * (EndTimeOfDay - StartTimeOfDay) / 630)
	StartX := (width / 25) + (int(event.Start.Weekday())-1)*eventWidth + 5
	canvas.Roundrect(StartX, StartY, eventWidth-10, eventHeight, 5, 5, "stroke:"+colorEvent(event.Category)+";stroke-width:3;fill:#32353B")

	canvas.Text(StartX+8, StartY+20, event.Module, "font-family:Roboto;fill:white")
	canvas.Text(StartX+8, StartY+38, event.Prof, "font-family:Roboto;fill:green")

	startTimeStr := fmt.Sprintf("%02dh%02d", event.Start.Hour(), event.Start.Minute())
	canvas.Text(StartX+eventWidth-len(startTimeStr)*13, StartY+22, startTimeStr, "font-family:Roboto;fill:white")

	endTimeStr := fmt.Sprintf("%02dh%02d", event.End.Hour(), event.End.Minute())
	canvas.Text(StartX+eventWidth-len(startTimeStr)*13, StartY+eventHeight-10, endTimeStr, "font-family:Roboto;fill:white")

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
	eventWidth := makeGrid(canvas, width, height)
	for _, event := range calendar {
		addEvent(canvas, width, height, eventWidth, event)
	}
	canvas.End()

	fileWriter.Flush()
	return s
}
