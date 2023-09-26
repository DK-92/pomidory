package main_view

import (
	"github.com/DK-92/pomidory/history"
	"github.com/DK-92/pomidory/view"
	"time"
)

func listenOnStateChannel() {
	for {
		value := <-stateChannel

		switch value {
		case view.PomodoroState:
			timerState = view.PomodoroState

			totalHistory.Add(pomodoroTimer.History, intentionInput.Text)
			pomodoroTimer.Stop()

			sleep()
			createInitialPomodoroView()

			pomodoroWindow.Show()
		case view.WorkBreakState:
			timerState = view.WorkBreakState

			pomodoroTimer.Stop()
			if history.GetInstance().IsBigBreak() {
				pomodoroTimer.Length = globalSettings.BigBreakLength
			} else {
				pomodoroTimer.Length = globalSettings.SmallBreakLength
			}

			sleep()

			startTimer()
			pomodoroWindow.Show()
		}
	}
}

func sleep() {
	time.Sleep(120 * time.Millisecond)
}
