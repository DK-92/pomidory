package main_view

import (
	"github.com/DK-92/pomidory/view"
	"time"
)

func listenOnStateChannel() {
	for {
		value := <-stateChannel

		switch value {
		case view.PomodoroState:
			timerState = view.PomodoroState
			pomodoroTimer.Stop()

			sleep()
			createInitialPomodoroView()

			pomodoroWindow.Show()
		case view.WorkBreakState:
			timerState = view.WorkBreakState

			pomodoroTimer.Stop()
			pomodoroTimer.Length = globalSettings.BreakLength
			sleep()

			startTimer()
			pomodoroWindow.Show()
		}
	}
}

func sleep() {
	time.Sleep(120 * time.Millisecond)
}
