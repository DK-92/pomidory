package work_break_view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/view"
)

func CreateAndShowWorkBreakView() {
	app := view.GetAppInstance()
	window := app.NewWindow("Work break")

	vbox := container.New(
		layout.NewVBoxLayout(),
		widget.NewRichTextFromMarkdown("## Work break!"),
		widget.NewLabel("Congratulations, 25 minutes has passed! It's now time for a quick 5 minute break."),
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
