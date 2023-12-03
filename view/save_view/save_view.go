package save_view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/DK-92/pomidory/view"
)

var (
	window fyne.Window
	isOpen bool
)

const (
	windowTitle = "Progress saved"

	saveText = `
Your progress has been saved! Check the folder which this program

resides in.
`
)

func CreateAndShowSaveView() {
	if isOpen {
		return
	}

	app := view.GetAppInstance()
	window = app.NewWindow(windowTitle)

	//totalHistory := history.GetInstance()
	//totalHistory.Save()

	vbox := container.New(
		layout.NewVBoxLayout(),
		widget.NewRichTextFromMarkdown(saveText),
		createButtonContainer(),
	)

	window.SetContent(vbox)
	window.SetFixedSize(true)

	isOpen = true
	window.Show()
}

func createButtonContainer() *fyne.Container {
	return container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		createCloseButton(),
		layout.NewSpacer(),
	)
}

func createCloseButton() *widget.Button {
	return widget.NewButton("Close", func() {
		isOpen = false
		window.Close()
	})
}
