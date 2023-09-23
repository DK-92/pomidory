package settings

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LightTheme = iota
	DarkTheme
)

const filename = "settings.json"

var (
	once     sync.Once
	instance *Settings
)

type Settings struct {
	PomodoroLength     time.Duration `json:"pomodoroLength"`
	BreakLength        time.Duration `json:"breakLength"`
	MinimizeAfterStart bool          `json:"minimizeAfterStart"`
	Theme              int8          `json:"theme"`
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
		Theme:              instance.Theme,
	}

	buffer, _ := json.Marshal(c)

	err := os.WriteFile(filename, buffer, 0644)
	if err != nil {
		log.Println("Error saving settings file: ", err)
	}
}

func (s *Settings) IsLightTheme() bool {
	return s.Theme == LightTheme
}

func loadSettings() *Settings {
	settings := &Settings{
		PomodoroLength:     25 * time.Minute,
		BreakLength:        5 * time.Minute,
		MinimizeAfterStart: true,
		Theme:              LightTheme,
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
