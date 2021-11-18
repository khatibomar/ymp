package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	ytExp = `http(?:s?):\/\/(?:www\.)?youtu(?:be\.com\/watch\?v=|\.be\/)([\w\-\_]*)(&(amp;)?[\w\?=]*)?`
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter youtube URL...")

	slider := widget.NewSlider(0, 100)

	linkValidator := validation.NewRegexp(ytExp, "Invalid youtube link")

	playBtn := widget.NewButton("play", func() {
		if input.Text == "" {
			err := fmt.Errorf("Empty link provided")
			log.Println(err)
			showError(myApp, err)
			return
		}
		if err := linkValidator(input.Text); err != nil {
			log.Println(err)
			showError(myApp, err)
			return
		}
		log.Println("Playing:", input.Text)
	})

	pauseBtn := widget.NewButton("pause", func() {
		log.Println("pause called")
	})

	content := container.New(
		layout.NewVBoxLayout(),
		input,
		playBtn,
		pauseBtn,
		slider,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func showError(a fyne.App, err error) {
	win := a.NewWindow("Error")
	win.Resize(fyne.NewSize(100, 50))
	content := container.NewVBox(
		widget.NewLabel(err.Error()),
		widget.NewButton("close", func() {
			win.Close()
		}),
	)
	win.SetContent(content)
	win.Show()
}
