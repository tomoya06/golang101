package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"syscall/js"

	"github.com/Lofanmi/chinese-calendar-golang/lunar"
	ics "github.com/arran4/golang-ical"
)

const DOWNLOAD_BASE_PATH = "./images/"
const MAX_LUNAR_AGES = 100
const LUNAR_LEAP_PREFIX = "(闰月)"

type LunarBirthday struct {
	YY   int    `json:"yy"`
	MM   int    `json:"mm"`
	DD   int    `json:"dd"`
	Name string `json:"name"`
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

func WriteToFile(input string) {
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
			yy:       birthday.YY + cnt,
			dd:       birthday.DD,
			mm:       birthday.MM,
			remindAt: remindAt,
			name:     birthday.Name,
		}
		bdInput.AddBirthdaysForOneYear(cal)
	}
}

func GenerateIcsContent(bds []LunarBirthday, remindAt int) string {
	cal := ics.NewCalendar()
	for _, bd := range bds {
		bd.AddBirthdays(cal, remindAt)
	}
	icsStr := cal.Serialize()
	return icsStr
}

type GenerateParams struct {
	BDS      []LunarBirthday `json:"bds"`
	RemindAt int             `json:"remindAt"`
}

func jsWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			fmt.Println("invalid args")
			return js.Null
		}
		fmt.Println("valid args")

		inputJSON := args[0].String()
		inputObj := GenerateParams{}
		json.Unmarshal([]byte(inputJSON), &inputObj)
		output := GenerateIcsContent(inputObj.BDS, inputObj.RemindAt)
		fmt.Println(inputJSON, "|", inputObj, "|", len(inputObj.BDS), inputObj.RemindAt)
		return output
	})
	return jsonFunc
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("generateIcsContent", jsWrapper())
	<-make(chan bool)
}

// func main() {
// 	inputJSON := `{"bds":[{"yy":1996,"mm":12,"dd":12,"name":"test"}],"remindAt":12}`
// 	inputObj := GenerateParams{}
// 	json.Unmarshal([]byte(inputJSON), &inputObj)
// 	output := GenerateIcsContent(inputObj.BDS, inputObj.RemindAt)
// 	fmt.Println(output)
// }
