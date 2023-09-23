package help_view

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
	windowTitle = "Help"

	helpText = `
This application is created to help you with the pomodoro technique. In the late 

1980's Francesco Cirillo developed this technique to help him work more efficiently.

This application will try and help you with that.

## How does the app work?
On the main screen, you'll notice an input field at the top, a counter below it  

and a button with the text 'Start Session'. Clicking the button will start a 

work session of 25 minutes. After that time is up, you'll receive a popup from 

the program urging you to take a quick 5 minute break. Make sure you don't 

do any work during this break! This 30 minute block is 1 pomodoro. Once you 

have completed 4 of them, it's time to take a longer 20 minute break.

After that, the process begins anew.

## Modifying timers
The timers' length can be modified by clicking the wheel icon in the top left 

corner of the application. The minimum time is 1 minute. You can adjust the 

work timer, as well as the break timer.

Furthermore, you can also enable or disable automatic closing of the main 

window after starting a new pomodoro. You may also change the theme from light 

to dark, if you so wish. 

Your preferences will be saved to a file called 'settings.json'.

## Author
This application is open source software under the GPL2 license. This excludes 

the logo, which is copyrighted to it's original author.
`
)

func CreateAndShowSettingsView() {
	if isOpen {
		return
	}

	app := view.GetAppInstance()
	window = app.NewWindow(windowTitle)

	vbox := container.New(
		layout.NewVBoxLayout(),
		widget.NewRichTextFromMarkdown(helpText),
		createOKButton(),
	)

	window.SetContent(vbox)
	window.Resize(fyne.Size{Height: 190, Width: 400})
	window.SetFixedSize(true)

	isOpen = true
	window.Show()
}

func createOKButton() *widget.Button {
	return widget.NewButton("Ok", func() {
		isOpen = false
		window.Close()
	})
}
