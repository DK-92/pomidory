package main_view

import (
	"github.com/DK-92/pomidory/view"
	"time"
)

func listenOnStateChannel() {
	for {
		value := <-channel

		switch value {
		case view.PomodoroState:
		case view.WorkBreakState:
			pomodoroTimer.Stop()
			pomodoroTimer.Length = 6 * time.Second
			pomodoroTimer.StartAfter(func() {
				time.Sleep(100 * time.Millisecond)
				createInitialPomodoroView()
			})
		}
	}
}
