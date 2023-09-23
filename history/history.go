package history

import (
	"fmt"
	"github.com/DK-92/pomidory/model"
	"log"
	"os"
	"sync"
	"time"
)

const filenamePrefix = "pomodoro_"

var (
	once     sync.Once
	instance *TotalHistory

	delimiter = ","
)

type TotalHistory struct {
	history     []*model.History
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

func (t *TotalHistory) Add(newHistory *model.History, task string) {
	newHistory.Task = task
	t.history = append(t.history, newHistory)
}

func (t *TotalHistory) Save() {
	filename := fmt.Sprintf("%s%s_%d.csv", filenamePrefix, time.Now().Format(model.DateFormat), t.currentTime)

	err := os.WriteFile(filename, []byte(t.toCSV()), 0644)
	if err != nil {
		log.Println("Error saving pomodoro history file: ", err)
	}
}

func (t *TotalHistory) toCSV() string {
	var output string

	output += "Date" + delimiter + "Length in HH MM" + delimiter + "Start Time" + delimiter + "End Time" + delimiter + "Task" + "\r\n"

	for _, oldHistory := range t.history {
		output += oldHistory.CurrentDate + delimiter
		output += oldHistory.Length + delimiter
		output += oldHistory.Start.Format(model.TimeFormat) + delimiter
		output += oldHistory.End.Format(model.TimeFormat) + delimiter
		output += oldHistory.Task + "\r\n"
	}

	return output
}
