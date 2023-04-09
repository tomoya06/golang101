package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Lofanmi/chinese-calendar-golang/lunar"
	ics "github.com/arran4/golang-ical"
)

const DOWNLOAD_BASE_PATH = "./images/"
const MAX_LUNAR_YEAR = 2100
const LUNAR_LEAP_PREFIX = "(闰月)"

type LunarEventBundle struct {
	tm   time.Time
	name string
}

func (evt *LunarEventBundle) AddEventToCalendar(cal *ics.Calendar) {
	tm, name := evt.tm, evt.name
	cal.SetMethod(ics.MethodRequest)
	evtId := fmt.Sprintf("lunarevent@%d#%s", tm.Unix(), strings.ReplaceAll(name, " ", "_"))
	event := cal.AddEvent(evtId)
	event.SetCreatedTime(tm)
	event.SetDtStampTime(tm)
	event.SetStartAt(tm)
	event.SetSummary(name)
}

func WriteToIcs(input string) {
	os.Mkdir(DOWNLOAD_BASE_PATH, os.ModePerm)

	f, err := os.Create("date.ics")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(input)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("done")
}

func main() {
	lunarYY, lunarMM, lunarDD := 1963, 4, 16
	evtName := "农历生日"
	cal := ics.NewCalendar()

	for ; lunarYY <= MAX_LUNAR_YEAR; lunarYY += 1 {
		solarTsNonLeap := lunar.ToSolarTimestamp(int64(lunarYY), int64(lunarMM), int64(lunarDD), 12, 0, 0, false)
		solarTsLeap := lunar.ToSolarTimestamp(int64(lunarYY), int64(lunarMM), int64(lunarDD), 12, 0, 0, true)

		tm := time.Unix(solarTsNonLeap, 0)
		evt := &LunarEventBundle{tm, evtName}
		evt.AddEventToCalendar(cal)

		if solarTsLeap != solarTsNonLeap {
			tm := time.Unix(solarTsNonLeap, 0)
			evt := &LunarEventBundle{tm, evtName + LUNAR_LEAP_PREFIX}
			evt.AddEventToCalendar(cal)
		}
	}

	icsStr := cal.Serialize()

	WriteToIcs(icsStr)
}
