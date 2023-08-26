package settings

import (
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *Settings
)

type Settings struct {
	PomodoroLength     time.Duration
	BreakLength        time.Duration
	MinimizeAfterStart bool
}

func GetInstance() *Settings {
	once.Do(func() {
		instance = &Settings{
			PomodoroLength:     5 * time.Second,
			BreakLength:        7 * time.Second,
			MinimizeAfterStart: false,
		}
	})

	return instance
}
