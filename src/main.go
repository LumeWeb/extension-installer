package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lumeweb/extension-installer/src/platform"
	_ "github.com/mediabuyerbot/go-crx3"
)

var appMain fyne.App

func main() {
	appMain = app.New()
	w := appMain.NewWindow("Lume Web Extension Installer")
	w.Resize(fyne.Size{Height: 250, Width: 250})

	content := widget.NewButton("Install", handleInstall)
	appMain.Settings().SetTheme(theme.LightTheme())
	w.SetContent(content)
	w.ShowAndRun()
}

func handleInstall() {
	go platform.StartInstall()
}
