package main_view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/model"
	"github.com/DK-92/pomidory/view"
	"github.com/DK-92/pomidory/view/work_break_view"
	"image/color"
	"time"
)

var (
	window fyne.Window

	timerText                 *canvas.Text
	intentionInput            *widget.Entry
	startTimerButtonContainer *fyne.Container
	stopTimerButtonContainer  *fyne.Container
	vbox                      *fyne.Container

	menu          *fyne.Menu
	menuRemainder *fyne.MenuItem

	pomodoroTimer *model.PomodoroTimer
)

func CreateAndShowMainView() {
	pomodoroTimer = model.GetInstance()
	app := app.New()
	window = app.NewWindow("Pomidory")
	createSystemTrayMenu()

	vbox = container.New(
		layout.NewVBoxLayout(),
		createIntentionInput(),
		createTimerText(pomodoroTimer.TimerLength()),
		createStartTimerButton(),
	)

	window.SetContent(vbox)
	window.Resize(fyne.NewSize(120, 120))
	window.SetFixedSize(true)
	window.CenterOnScreen()

	window.SetCloseIntercept(func() {
		window.Hide()
	})

	window.ShowAndRun()
}

func createSystemTrayMenu() {
	app := view.GetAppInstance()
	menuRemainder = fyne.NewMenuItem("Timer not started", nil)

	if desk, isDesktop := app.(desktop.App); isDesktop {
		menu = fyne.NewMenu("Pomidory",
			fyne.NewMenuItem("Show", func() {
				window.Show()
			}),
			fyne.NewMenuItemSeparator(),
			menuRemainder,
		)
		desk.SetSystemTrayMenu(menu)
	}
}

func createIntentionInput() *widget.Entry {
	intentionInput = widget.NewEntry()
	intentionInput.SetPlaceHolder("Task at hand")

	return intentionInput
}

func createStartTimerButton() *fyne.Container {
	if startTimerButtonContainer != nil {
		return startTimerButtonContainer
	}

	startTimerButton := widget.NewButton("Start session", func() {
		pomodoroTimer.StartAfter(func() {
			work_break_view.CreateAndShowWorkBreakView()
		})

		// Remove the 3rd item from layout (start timer button)
		vbox.Objects = vbox.Objects[:2]
		vbox.Add(createStopTimerButton())

		intentionInput.Disable()
		updateTimerText(pomodoroTimer.Remainder())

		// Update the time element on the UI
		go func() {
			for range time.Tick(970 * time.Millisecond) {
				if pomodoroTimer.HasEnded() {
					return
				}

				remainder := pomodoroTimer.Remainder()
				updateTimerText(remainder)
				updateMenuItemTimerText(fmt.Sprintf("Time left: %s", remainder))
			}
		}()

		// Close after 2 seconds, so the user sees the timer has started
		//go func() {
		//	time.Sleep(2 * time.Second)
		//	window.Hide()
		//}()
	})

	startTimerButtonContainer = container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		startTimerButton,
		layout.NewSpacer(),
	)

	return startTimerButtonContainer
}

func createStopTimerButton() *fyne.Container {
	if stopTimerButtonContainer != nil {
		return stopTimerButtonContainer
	}

	stopTimerButton := widget.NewButton("Stop session", func() {
		pomodoroTimer.Stop()
		intentionInput.Enable()

		vbox.Objects = vbox.Objects[:2]
		vbox.Add(startTimerButtonContainer)
	})

	stopTimerButtonContainer = container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		stopTimerButton,
		layout.NewSpacer(),
	)

	return stopTimerButtonContainer
}

func createTimerText(text string) *canvas.Text {
	timerText = canvas.NewText(text, color.Black)
	timerText.TextSize = 40
	timerText.Alignment = fyne.TextAlignCenter
	return timerText
}

func updateTimerText(text string) {
	timerText.Text = text
	vbox.Refresh()
}

func updateMenuItemTimerText(text string) {
	menuRemainder.Label = text
	menu.Refresh()
}
