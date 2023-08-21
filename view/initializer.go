package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"sync"
)

var (
	once        sync.Once
	application fyne.App
)

func GetAppInstance() fyne.App {
	once.Do(func() {
		application = app.New()
	})

	return application
}
