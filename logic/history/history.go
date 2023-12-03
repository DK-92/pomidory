package history

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	filenamePrefix      = "pomodoro_"
	bigWorkBreakCounter = 4

	dateFormat = "01-02-2006"
	timeFormat = "15:04"
)

var (
	once     sync.Once
	instance *TotalHistory

	delimiter = ","
)

type History struct {
	CurrentDate string
	Start       time.Time
	End         time.Time
	Task        string
}

func NewHistory(now time.Time, length time.Duration, task string) *History {
	return &History{
		CurrentDate: now.Format(dateFormat),
		Start:       now,
		End:         now.Add(length),
		Task:        task,
	}
}

type TotalHistory struct {
	history     []*History
	currentTime int64
}

func GetInstance() *TotalHistory {
	once.Do(func() {
		instance = &TotalHistory{
			currentTime: time.Now().Unix(),
		}
	})

	return instance
}

func (t *TotalHistory) Add(newHistory *History) {
	t.history = append(t.history, newHistory)
}

func (t *TotalHistory) Save() {
	filename := fmt.Sprintf("%s%s_%d.csv", filenamePrefix, time.Now().Format(dateFormat), t.currentTime)

	err := os.WriteFile(filename, []byte(t.toCSV()), 0644)
	if err != nil {
		log.Println("Error saving pomodoro history file: ", err)
	}
}

func (t *TotalHistory) IsBigBreak() bool {
	println(len(t.history))

	if len(t.history) == 0 {
		return false
	}

	return ((len(t.history) + 1) % bigWorkBreakCounter) == 0
}

func (t *TotalHistory) toCSV() string {
	var output string

	output += "Date" + delimiter + "Start Time" + delimiter + "End Time" + delimiter + "Task" + "\r\n"

	for _, oldHistory := range t.history {
		output += oldHistory.CurrentDate + delimiter
		output += oldHistory.Start.Format(timeFormat) + delimiter
		output += oldHistory.End.Format(timeFormat) + delimiter
		output += oldHistory.Task + "\r\n"
	}

	return output
}
