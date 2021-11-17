package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("ymp - for gamers")
	w.SetContent(widget.NewLabel("Hello Fyne!"))
	w.SetContent(widget.NewButton("play", func() {
		log.Println("Button pressed")
	}))

	w.ShowAndRun()
}
