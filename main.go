package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	slider := widget.NewSlider(0, 100)

	playBtn := widget.NewButton("play", func() {
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
