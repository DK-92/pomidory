package main_view

import (
	"github.com/DK-92/pomidory/view"
	"time"
)

func listenOnStateChannel() {
	for {
		value := <-channel

		switch value {
		case view.WorkBreakState:
			pomodoroTimer.Stop()
			pomodoroTimer.Length = globalSettings.BreakLength
			pomodoroWindow.Show()

			pomodoroTimer.StartAfter(func() {
				time.Sleep(100 * time.Millisecond)
				createInitialPomodoroView()
			})
		}
	}
}
