package main_view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/view/work_break_view"
	"time"
)

const hideWindowAfterStartTimerSeconds = 2 * time.Second

func createInitialPomodoroView() {
	if vbox == nil {
		vbox = container.New(
			layout.NewVBoxLayout(),
			createOrSetIntentionInput(),
			createOrUpdateTimerText(pomodoroTimer.TimerLength()),
			createOrSetStartTimerButton(),
		)
	}

	pomodoroTimer.Length = globalSettings.PomodoroLength
	createOrUpdateTimerText(pomodoroTimer.TimerLength())

	intentionInput.Text = ""
	intentionInput.Enable()
	vbox.Refresh()
}

func createOrSetIntentionInput() *widget.Entry {
	if intentionInput == nil {
		intentionInput = widget.NewEntry()
	}

	intentionInput.SetPlaceHolder("Task at hand")
	return intentionInput
}

func createOrSetStartTimerButton() *fyne.Container {
	if startTimerButtonContainer != nil {
		return startTimerButtonContainer
	}

	startTimerButton := widget.NewButton("Start session", func() {
		pomodoroTimer.StartAfter(func() {
			work_break_view.CreateAndShowWorkBreakView(channel)
		})

		// Remove the 3rd item from layout (start timer button)
		vbox.Objects = vbox.Objects[:2]
		vbox.Add(createOrSetStopTimerButton())

		intentionInput.Disable()
		createOrUpdateTimerText(pomodoroTimer.Remainder())

		// Update the time element on the UI
		go func() {
			for range time.Tick(60 * time.Millisecond) {
				if pomodoroTimer.HasEnded() {
					addStartButtonToContainer()

					return
				}

				remainder := pomodoroTimer.Remainder()
				createOrUpdateTimerText(remainder)
				updateMenuItemTimerText(fmt.Sprintf("Time left: %s", remainder))
			}
		}()

		// Close after 2 seconds, so the user sees the timer has started
		if globalSettings.MinimizeAfterStart {
			go func() {
				time.Sleep(hideWindowAfterStartTimerSeconds)
				pomodoroWindow.Hide()
			}()
		}
	})

	startTimerButtonContainer = container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		startTimerButton,
		layout.NewSpacer(),
	)

	return startTimerButtonContainer
}

func createOrSetStopTimerButton() *fyne.Container {
	if stopTimerButtonContainer != nil {
		return stopTimerButtonContainer
	}

	stopTimerButton := widget.NewButton("Stop session", func() {
		pomodoroTimer.Stop()
		intentionInput.Enable()
		addStartButtonToContainer()
	})

	stopTimerButtonContainer = container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		stopTimerButton,
		layout.NewSpacer(),
	)

	return stopTimerButtonContainer
}

func addStartButtonToContainer() {
	vbox.Objects = vbox.Objects[:2]
	vbox.Add(startTimerButtonContainer)
}
