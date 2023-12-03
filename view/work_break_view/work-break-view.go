package work_break_view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/settings"
	"github.com/DK-92/pomidory/view"
	"time"
)

const (
	windowTitle = "Work break"

	ShortBreak = iota
	LongBreak
)

var (
	globalSettings *settings.Settings

	breakType int
)

func CreateAndShowWorkBreakView(bType int) {
	app := view.GetAppInstance()
	window := app.NewWindow(windowTitle)

	globalSettings = settings.GetInstance()
	breakType = bType

	vbox := container.New(
		layout.NewVBoxLayout(),
		widget.NewRichTextFromMarkdown(createLabelText()),
		createButtons(window),
	)

	window.SetContent(vbox)
	window.SetFixedSize(true)
	window.CenterOnScreen()
	window.Show()
}

func createButtons(window fyne.Window) *fyne.Container {
	button := widget.NewButton("OK", func() {
		window.Close()
	})

	hbox := container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		button,
		layout.NewSpacer(),
	)

	return hbox
}

func createLabelText() string {
	var result string

	result = fmt.Sprintf(`
## Work break!

Congratulations, %d minutes has passed!

%s
`, globalSettings.PomodoroLength/time.Minute, createBreakText())

	return result
}

func createBreakText() string {
	if breakType == LongBreak {
		return fmt.Sprintf("It's time for a longer %d minute break.", globalSettings.BigBreakLength/time.Minute)
	}

	return fmt.Sprintf("It's time for a quick %d minute break.", globalSettings.SmallBreakLength/time.Minute)
}
