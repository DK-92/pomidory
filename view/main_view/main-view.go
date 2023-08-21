package main_view

import (
	"fyne.io/fyne/v2"
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
	timerText *canvas.Text
	vbox      *fyne.Container

	menu          *fyne.Menu
	menuRemainder *fyne.MenuItem

	pomodoroTimer = model.GetInstance()
)

func CreateAndShowMainView() {
	app := view.GetAppInstance()
	window := app.NewWindow("Pomidory")
	createSystemTrayMenu(window)

	vbox = container.New(
		layout.NewVBoxLayout(),
		createIntentionInput(),
		createTimerText("25:00"),
		createStartTimerButton(window),
	)

	window.SetContent(vbox)
	window.Resize(fyne.NewSize(120, 120))
	window.SetFixedSize(true)
	window.CenterOnScreen()
	window.ShowAndRun()

	window.SetCloseIntercept(func() {
		window.Hide()
	})
}

func createSystemTrayMenu(window fyne.Window) {
	app := view.GetAppInstance()

	menuRemainder = fyne.NewMenuItem("Time left: 25:00", func() {})

	if desk, isDesktop := app.(desktop.App); isDesktop {
		menu = fyne.NewMenu("Pomidory",
			fyne.NewMenuItem("Show", func() {
				println("Opening main window")
				window.Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Time left: 25:00", func() {}),
		)
		desk.SetSystemTrayMenu(menu)
	}
}

func createIntentionInput() *widget.Entry {
	input := widget.NewEntry()
	input.SetPlaceHolder("Task at hand")

	return input
}

func createStartTimerButton(window fyne.Window) *fyne.Container {
	startTimerButton := widget.NewButton("Start session", func() {

		updateTimerText(pomodoroTimer.Remainder())

		go func() {
			for pomodoroTimer.HasNotEnded() {
				remainder := pomodoroTimer.Remainder()
				updateTimerText(remainder)
				go updateMenuItemTimerText(remainder)
				time.Sleep(980 * time.Millisecond)
			}
		}()

		pomodoroTimer.Start(func() {
			work_break_view.CreateAndShowWorkBreakView()
		})

		window.Hide()
	})

	return container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		startTimerButton,
		layout.NewSpacer(),
	)
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
