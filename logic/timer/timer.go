package timer

import (
	"fmt"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *Timer
)

const (
	Stopped = iota
	Pomodoro
	Break
)

type Timer struct {
	Length  time.Duration
	running *time.Timer
	start   time.Time
	end     time.Time
}

func GetInstance() *Timer {
	once.Do(func() {
		instance = &Timer{
			Length: 2 * time.Second,
		}
	})

	return instance
}

func (t *Timer) StartAndRunAfter(runAfter func()) {
	t.running = time.AfterFunc(t.Length, func() {
		runAfter()
	})

	t.start = time.Now()
	t.end = time.Now().Add(t.Length)
}

func (t *Timer) TimerLength() string {
	t.end = time.Now().Add(t.Length)
	return t.Remainder()
}

func (t *Timer) Remainder() string {
	difference := t.end.Sub(time.Now()).Round(time.Second)

	total := int(difference.Seconds())
	minutes := int(total/60) % 60
	seconds := total % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (t *Timer) HasEnded() bool {
	return time.Now().After(t.end)
}

func (t *Timer) Stop() {
	t.end = time.UnixMilli(0)
	t.running.Stop()
}
