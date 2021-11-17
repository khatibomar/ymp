package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

const (
	ytExp = `http(?:s?):\/\/(?:www\.)?youtu(?:be\.com\/watch\?v=|\.be\/)([\w\-\_]*)(&(amp;)?[\w\?=]*)?`
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	slider := widget.NewSlider(0, 100)

	linkValidator := validation.NewRegexp(ytExp, "Invalid youtube link")

	playBtn := widget.NewButton("play", func() {
		if input.Text == "" {
			log.Println("Empty link provided")
			return
		}
		if err := linkValidator(input.Text); err != nil {
			log.Println(err)
			return
		}
		log.Println("Playing:", input.Text)
	})

	pauseBtn := widget.NewButton("pause", func() {
		log.Println("pause called")
	})

	content := container.NewVBox(
		input,
		container.NewHBox(
			playBtn,
			pauseBtn,
			// TODO(khatibomar) : Investigate why slider is not taking all space
			container.NewMax(
				slider,
			),
		),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
