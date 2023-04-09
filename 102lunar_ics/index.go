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
const MAX_LUNAR_AGES = 100
const LUNAR_LEAP_PREFIX = "(闰月)"

type LunarBirthday struct {
	yy   int
	mm   int
	dd   int
	name string
}

type LunarEventInputs struct {
	yy       int
	dd       int
	mm       int
	remindAt int
	name     string
}

type LunarEventBundle struct {
	tm   time.Time
	name string
}

func (input *LunarEventInputs) AddBirthdaysForOneYear(cal *ics.Calendar) {
	lunarYY, lunarMM, lunarDD, evtName := input.yy, input.mm, input.dd, input.name
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

	f, err := os.Create(DOWNLOAD_BASE_PATH + "date.ics")
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

func (birthday *LunarBirthday) AddBirthdays(cal *ics.Calendar, remindAt int) {
	for cnt := 0; cnt <= MAX_LUNAR_AGES; cnt += 1 {
		bdInput := LunarEventInputs{
			yy:       birthday.yy + cnt,
			dd:       birthday.dd,
			mm:       birthday.mm,
			remindAt: remindAt,
			name:     birthday.name,
		}
		bdInput.AddBirthdaysForOneYear(cal)
	}
}

func main() {
	bds := []LunarBirthday{
		{
			yy:   1963,
			mm:   4,
			dd:   16,
			name: "🎂父亲大人生日🎂",
		},
		{
			yy:   1967,
			mm:   9,
			dd:   23,
			name: "🎂母亲大人生日🎂",
		},
	}
	remindAt := 12

	cal := ics.NewCalendar()
	for _, bd := range bds {
		bd.AddBirthdays(cal, remindAt)
	}
	icsStr := cal.Serialize()
	WriteToIcs(icsStr)
}
