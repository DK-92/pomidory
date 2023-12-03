package main_view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/logic/history"
	"github.com/DK-92/pomidory/logic/timer"
	"github.com/DK-92/pomidory/settings"
	"github.com/DK-92/pomidory/view"
	"github.com/DK-92/pomidory/view/help_view"
	"github.com/DK-92/pomidory/view/save_view"
	"github.com/DK-92/pomidory/view/settings_view"
	"github.com/DK-92/pomidory/view/work_break_view"
	"log/slog"
	"sync/atomic"
	"time"
)

const (
	windowTitle = "Pomidory"

	hideWindowAfterStartTimerSeconds = 2 * time.Second
	buttonPositionInVbox             = 3
	bigBreak                         = 4
)

const (
	stateStopped = iota
	statePomodoro
	stateBreak
	stateBreakFinished
)

var (
	pomodoroWindow fyne.Window

	timerText                 *canvas.Text
	intentionInput            *widget.Entry
	startTimerButtonContainer *fyne.Container
	stopTimerButtonContainer  *fyne.Container
	vbox                      *fyne.Container

	menu          *fyne.Menu
	menuRemainder *fyne.MenuItem

	globalSettings *settings.Settings
	globalHistory  *history.TotalHistory
	pTimer         *timer.Timer

	pomodorosFinished atomic.Int32
	state             atomic.Int32
)

func CreateAndShowMainView() {
	globalSettings = settings.GetInstance()
	globalHistory = history.GetInstance()
	pTimer = timer.GetInstance()

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

func createInitialPomodoroView() {
	if vbox == nil {
		vbox = container.New(
			layout.NewVBoxLayout(),
			createToolbar(),
			createOrSetIntentionInput(),
			createOrUpdateTimerText(""),
			createOrSetStartTimerButton(),
		)
	}

	// TODO: Uncomment below after work is done
	pTimer.Length = globalSettings.PomodoroLength
	createOrUpdateTimerText(pTimer.TimerLength())

	intentionInput.Text = ""
	intentionInput.Enable()
	vbox.Refresh()
}

func createToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), settings_view.CreateAndShowSettingsView),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), save_view.CreateAndShowSaveView),
		widget.NewToolbarAction(theme.HelpIcon(), help_view.CreateAndShowHelpView),
	)
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

	startTimerButton := widget.NewButton("Start session", startTimer)

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
		pTimer.Stop()
		pTimer.Length = globalSettings.PomodoroLength
		state.Store(stateStopped)

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
	vbox.Objects = vbox.Objects[:buttonPositionInVbox]
	vbox.Add(startTimerButtonContainer)
	pTimer.Length = globalSettings.PomodoroLength
}

func startTimer() {
	now := time.Now()

	if state.Load() == stateStopped {
		state.Store(statePomodoro)
		globalHistory.Add(history.NewHistory(now, pTimer.Length, intentionInput.Text))
		slog.Info(windowTitle, "message", "Starting pomodoro", "state", state.Load())
	}

	pTimer.StartAndRunAfter(now, func() {
		time.Sleep(120 * time.Millisecond)

		if state.Load() == statePomodoro {
			pomodorosFinished.Add(1)
			slog.Info(windowTitle, "message", "Pomodoro finished", "finishedPomodoros", pomodorosFinished.Load(), "state", state.Load())
			intentionInput.Enable()
			addStartButtonToContainer()
		}

		if state.Load() == stateBreakFinished {
			state.Store(stateStopped)
			slog.Info(windowTitle, "message", "Resetting state", "state", state.Load())
			intentionInput.Enable()
			addStartButtonToContainer()
			return
		}

		if pomodorosFinished.Load()%bigBreak == 0 {
			slog.Info(windowTitle, "message", "Spawing work break view", "breakType", "LongBreak")
			work_break_view.CreateAndShowWorkBreakView(work_break_view.LongBreak)
			pomodorosFinished.Store(0)

			pTimer.Length = globalSettings.BigBreakLength
			state.Store(stateBreak)
		} else {
			slog.Info(windowTitle, "message", "Spawing work break view", "breakType", "ShortBreak")
			work_break_view.CreateAndShowWorkBreakView(work_break_view.ShortBreak)

			pTimer.Length = globalSettings.SmallBreakLength
			state.Store(stateBreak)
		}

		if state.Load() == stateBreak {
			slog.Info(windowTitle, "message", "Starting work break timer", "state", state.Load())
			startTimer()
		}
	})

	// Remove the 3rd item from layout (start timer button)
	vbox.Objects = vbox.Objects[:buttonPositionInVbox]
	vbox.Add(createOrSetStopTimerButton())

	intentionInput.Disable()
	createOrUpdateTimerText(pTimer.Remainder())

	// Update the time element on the UI
	go func() {
		for range time.Tick(60 * time.Millisecond) {
			if pTimer.HasEnded() {
				if state.Load() == stateBreak {
					state.Store(stateBreakFinished)
				}

				slog.Info(windowTitle, "message", "Timer text loop stopped", "state", state.Load())
				return
			}

			remainder := pTimer.Remainder()
			createOrUpdateTimerText(remainder)
			updateMenuItemTimerText(fmt.Sprintf("Time left: %s", remainder))
		}
	}()

	// Close after 2 seconds, so the user sees the timer has started
	if globalSettings.MinimizeAfterStart {
		go func() {
			time.Sleep(hideWindowAfterStartTimerSeconds)
			pomodoroWindow.Hide()
			slog.Info(windowTitle, "message", "MinimizeAfterStart")
		}()
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

	if text != "" {
		timerText.Text = text
		timerText.Refresh()
	}

	return timerText
}

func updateMenuItemTimerText(text string) {
	menuRemainder.Label = text
	menu.Refresh()
}
