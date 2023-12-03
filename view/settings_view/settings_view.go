package settings_view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/components"
	"github.com/DK-92/pomidory/settings"
	"github.com/DK-92/pomidory/view"
	"time"
)

const (
	windowTitle = "Settings"

	lightTheme = "Light"
	darkTheme  = "Dark"
)

var (
	//pomodoroTimer  *pomodoro.PomodoroTimer
	globalSettings *settings.Settings

	window fyne.Window
	isOpen bool

	pomodoroTimerLength   *components.NumericalEntry
	smallBreakTimerLength *components.NumericalEntry
	bigBreakTimerLength   *components.NumericalEntry

	windowClose     *widget.Check
	themeRadioGroup *widget.RadioGroup
)

func CreateAndShowSettingsView() {
	if isOpen {
		return
	}

	//pomodoroTimer = pomodoro.GetInstance()
	globalSettings = settings.GetInstance()

	app := view.GetAppInstance()
	window = app.NewWindow(windowTitle)

	vbox := container.New(
		layout.NewVBoxLayout(),
		widget.NewRichTextFromMarkdown("All times below are given in full minutes.\n\nThe minimum time cannot be less than one minute."),
		createSettingsForm(),
		createButtons(),
	)

	window.SetContent(vbox)
	window.Resize(fyne.Size{Height: 190, Width: 400})
	window.SetFixedSize(true)

	isOpen = true
	window.Show()
}

func createSettingsForm() *fyne.Container {
	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Pomodoro length"),
		createPomodoroTimerLengthEntry(),
		widget.NewLabel("Small break length"),
		createSmallBreakTimerLengthEntry(),
		widget.NewLabel("Big break length"),
		createBigBreakTimerLengthEntry(),
		widget.NewLabel("Window closing"),
		createCloseWindowAutomaticallyCheck(),
		widget.NewLabel("App theme"),
		createThemeRadioButtons(),
	)
}

func createPomodoroTimerLengthEntry() *components.NumericalEntry {
	pomodoroTimerLength = components.NewNumericalEntry()
	pomodoroTimerLength.SetText(fmt.Sprintf("%.0f", globalSettings.PomodoroLength.Minutes()))

	return pomodoroTimerLength
}

func createSmallBreakTimerLengthEntry() *components.NumericalEntry {
	smallBreakTimerLength = components.NewNumericalEntry()
	smallBreakTimerLength.SetText(fmt.Sprintf("%.0f", globalSettings.SmallBreakLength.Minutes()))

	return smallBreakTimerLength
}

func createBigBreakTimerLengthEntry() *components.NumericalEntry {
	bigBreakTimerLength = components.NewNumericalEntry()
	bigBreakTimerLength.SetText(fmt.Sprintf("%.0f", globalSettings.BigBreakLength.Minutes()))

	return bigBreakTimerLength
}

func createCloseWindowAutomaticallyCheck() *widget.Check {
	windowClose = widget.NewCheck("", nil)
	windowClose.Checked = globalSettings.MinimizeAfterStart
	return windowClose
}

func createThemeRadioButtons() *widget.RadioGroup {
	themeRadioGroup = widget.NewRadioGroup([]string{lightTheme, darkTheme}, nil)

	if globalSettings.IsLightTheme() {
		themeRadioGroup.Selected = lightTheme
	} else {
		themeRadioGroup.Selected = darkTheme
	}

	return themeRadioGroup
}

func createButtons() *fyne.Container {
	return container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		createCancelButton(),
		createSaveButton(),
	)
}

func createCancelButton() *widget.Button {
	return widget.NewButton("Cancel", func() {
		isOpen = false
		window.Close()
	})
}

func createSaveButton() *widget.Button {
	return widget.NewButton("Save", func() {
		pomodoroTimerDuration, _ := time.ParseDuration(pomodoroTimerLength.Text + "m")
		globalSettings.PomodoroLength = pomodoroTimerDuration
		//pomodoroTimer.Length = globalSettings.PomodoroLength

		smallBreakTimerDuration, _ := time.ParseDuration(smallBreakTimerLength.Text + "m")
		globalSettings.SmallBreakLength = smallBreakTimerDuration

		bigBreakTimerDuration, _ := time.ParseDuration(bigBreakTimerLength.Text + "m")
		globalSettings.BigBreakLength = bigBreakTimerDuration

		globalSettings.MinimizeAfterStart = windowClose.Checked

		app := view.GetAppInstance()
		if themeRadioGroup.Selected == lightTheme {
			globalSettings.Theme = settings.LightTheme
			app.Settings().SetTheme(theme.LightTheme())
		} else {
			globalSettings.Theme = settings.DarkTheme
			app.Settings().SetTheme(theme.DarkTheme())
		}

		globalSettings.Save()
		isOpen = false
		window.Close()
	})
}
