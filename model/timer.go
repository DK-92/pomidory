package model

import (
	"fmt"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *PomodoroTimer
)

type PomodoroTimer struct {
	Length  int64
	running *time.Timer
	start   time.Time
	end     time.Time
}

func GetInstance() *PomodoroTimer {
	once.Do(func() {
		instance = &PomodoroTimer{Length: 1}
	})

	return instance
}

func (t *PomodoroTimer) Start(runAfter func()) {
	t.running = time.AfterFunc(time.Duration(t.Length)*time.Minute, runAfter)
	t.start = time.Now()
	t.end = time.Now().Add(time.Duration(t.Length) * time.Minute)
}

func (t *PomodoroTimer) Remainder() string {
	difference := t.end.Sub(time.Now())

	total := int(difference.Seconds())
	minutes := int(total/60) % 60
	seconds := total % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (t *PomodoroTimer) HasNotEnded() bool {
	return t.end.After(time.Now())
}
