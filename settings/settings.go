package settings

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"runtime"
	"sync"
	"time"
)

const (
	OSDefault = iota
	LightTheme
	DarkTheme
)

const (
	filename = ".pomidory_settings.json"

	defaultPomodoroLength   = 25 * time.Minute
	defaultSmallBreakLength = 5 * time.Minute
	defaultBigBreakLength   = 20 * time.Minute
)

var (
	once     sync.Once
	instance *Settings
)

type Settings struct {
	PomodoroLength     time.Duration `json:"pomodoroLength"`
	SmallBreakLength   time.Duration `json:"smallBreakLength"`
	BigBreakLength     time.Duration `json:"bigBreakLength"`
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
		SmallBreakLength:   instance.SmallBreakLength / time.Minute,
		BigBreakLength:     instance.BigBreakLength / time.Minute,
		MinimizeAfterStart: instance.MinimizeAfterStart,
		Theme:              instance.Theme,
	}

	buffer, _ := json.Marshal(c)

	err := os.WriteFile(fetchHomeDirIfNotWindows()+filename, buffer, 0644)
	if err != nil {
		log.Println("Error saving settings file: ", err)
	}
}

func loadSettings() *Settings {
	settings := &Settings{
		PomodoroLength:     25 * time.Minute,
		SmallBreakLength:   5 * time.Minute,
		BigBreakLength:     20 * time.Minute,
		MinimizeAfterStart: true,
		Theme:              OSDefault,
	}

	buffer, err := os.ReadFile(fetchHomeDirIfNotWindows() + filename)
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
	settings.SmallBreakLength = settings.SmallBreakLength * time.Minute
	settings.BigBreakLength = settings.BigBreakLength * time.Minute

	checkProperValues(settings)

	return settings
}

func checkProperValues(settings *Settings) {
	if isOutOfRange(settings.PomodoroLength) {
		settings.PomodoroLength = defaultPomodoroLength
	}

	if isOutOfRange(settings.SmallBreakLength) {
		settings.SmallBreakLength = defaultSmallBreakLength
	}

	if isOutOfRange(settings.BigBreakLength) {
		settings.BigBreakLength = defaultBigBreakLength
	}
}

func isOutOfRange(value time.Duration) bool {
	return value.Minutes() < 1 || value.Minutes() > 999
}

func fetchHomeDirIfNotWindows() string {
	if runtime.GOOS != "darwin" {
		return ""
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Println("Error fetching current user: ", err)
	}

	return currentUser.HomeDir + "/"
}
