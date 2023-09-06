package settings

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *Settings
)

const filename = "settings.json"

type Settings struct {
	PomodoroLength     time.Duration `json:"pomodoroLength"`
	BreakLength        time.Duration `json:"breakLength"`
	MinimizeAfterStart bool          `json:"minimizeAfterStart"`
}

func GetInstance() *Settings {
	once.Do(func() {
		instance = loadSettings()
	})

	return instance
}

func (s *Settings) Save() {
	c := &Settings{
		PomodoroLength:     instance.PomodoroLength / time.Minute,
		BreakLength:        instance.BreakLength / time.Minute,
		MinimizeAfterStart: instance.MinimizeAfterStart,
	}

	buffer, _ := json.Marshal(c)

	err := os.WriteFile(filename, buffer, 0644)
	if err != nil {
		log.Println("Error saving settings file: ", err)
	}
}

func loadSettings() *Settings {
	settings := &Settings{
		PomodoroLength:     25 * time.Minute,
		BreakLength:        5 * time.Minute,
		MinimizeAfterStart: true,
	}

	buffer, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error opening settings file: ", err)
		return settings
	}

	err = json.Unmarshal(buffer, &settings)
	if err != nil {
		log.Println("Error unmarshaling file: ", err)
		return settings
	}

	settings.PomodoroLength = settings.PomodoroLength * time.Minute
	settings.BreakLength = settings.BreakLength * time.Minute

	return settings
}
