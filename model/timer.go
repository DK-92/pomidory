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

const LengthType = time.Second

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

func (t *PomodoroTimer) TimerLength() string {
	t.end = time.Now().Add(time.Duration(t.Length) * LengthType)
	return t.Remainder()
}

func (t *PomodoroTimer) StartAfter(runAfter func()) {
	t.running = time.AfterFunc(time.Duration(t.Length)*LengthType, runAfter)
	t.start = time.Now()
	t.end = time.Now().Add(time.Duration(t.Length) * LengthType)
}

func (t *PomodoroTimer) Remainder() string {
	difference := t.end.Sub(time.Now()).Round(time.Second)

	total := int(difference.Seconds())
	minutes := int(total/60) % 60
	seconds := total % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (t *PomodoroTimer) HasEnded() bool {
	return time.Now().After(t.end)
}

func (t *PomodoroTimer) Stop() {
	t.end = time.UnixMilli(0)
	t.running.Stop()
}
