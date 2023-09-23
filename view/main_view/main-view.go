package main_view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/history"
	"github.com/DK-92/pomidory/pomodoro"
	"github.com/DK-92/pomidory/settings"
	"github.com/DK-92/pomidory/view"
)

const windowTitle = "Pomidory"

var (
	pomodoroWindow fyne.Window

	timerText                 *canvas.Text
	intentionInput            *widget.Entry
	startTimerButtonContainer *fyne.Container
	stopTimerButtonContainer  *fyne.Container
	vbox                      *fyne.Container

	menu          *fyne.Menu
	menuRemainder *fyne.MenuItem

	pomodoroTimer  *pomodoro.PomodoroTimer
	globalSettings *settings.Settings
	totalHistory   *history.TotalHistory

	stateChannel chan view.StateChannel
)

func CreateAndShowMainView() {
	pomodoroTimer = pomodoro.GetInstance()
	globalSettings = settings.GetInstance()
	totalHistory = history.GetInstance()

	stateChannel = make(chan view.StateChannel)
	go listenOnStateChannel()

	app := view.GetAppInstance()
	pomodoroWindow = app.NewWindow(windowTitle)
	createSystemTrayMenu()
	createInitialPomodoroView()

	pomodoroWindow.SetContent(vbox)
	pomodoroWindow.Resize(fyne.NewSize(120, 120))
	pomodoroWindow.SetFixedSize(true)
	pomodoroWindow.CenterOnScreen()

	pomodoroWindow.SetCloseIntercept(func() {
		pomodoroWindow.Hide()
	})

	setInitialTheme()
	pomodoroWindow.ShowAndRun()
}

func createSystemTrayMenu() {
	app := view.GetAppInstance()
	menuRemainder = fyne.NewMenuItem("Timer not started", nil)

	if desk, isDesktop := app.(desktop.App); isDesktop {
		menu = fyne.NewMenu("Pomidory",
			fyne.NewMenuItem("Show", func() {
				pomodoroWindow.Show()
			}),
			fyne.NewMenuItemSeparator(),
			menuRemainder,
		)
		desk.SetSystemTrayMenu(menu)
	}
}

func setInitialTheme() {
	app := view.GetAppInstance()
	// TODO: This will have to be refactored in fyne v3
	if globalSettings.IsLightTheme() {
		app.Settings().SetTheme(theme.LightTheme())
	} else {
		app.Settings().SetTheme(theme.DarkTheme())
	}
}

func createOrUpdateTimerText(text string) *canvas.Text {
	if timerText == nil {
		// Interesting case, setting color nil will automatically update color on theme change
		// Worth keeping in mind when updating fyne
		timerText = canvas.NewText(text, nil)
		timerText.Alignment = fyne.TextAlignCenter
		timerText.TextSize = 40
	}
	timerText.Text = text
	timerText.Refresh()
	return timerText
}

func updateMenuItemTimerText(text string) {
	menuRemainder.Label = text
	menu.Refresh()
}
