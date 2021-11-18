package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	remoteMpv "github.com/blang/mpv"
	"github.com/khatibomar/ymp/mpv"
	"github.com/pkg/errors"
)

const (
	ytExp = `http(?:s?):\/\/(?:www\.)?youtu(?:be\.com\/watch\?v=|\.be\/)([\w\-\_]*)(&(amp;)?[\w\?=]*)?`
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ymp - for gamers")
	log.Println("Starting application...")
	mpv, err, cleanup := mpv.New()
	if err != nil {
		log.Println(err)
		showError(myApp, err)
	}
	defer cleanup()
	log.Println("Creating mpv instance")

	if err := mpv.Cmd.Start(); err != nil {
		log.Println(errors.Wrap(err, "failed to start mpv"))
	}
	if err := mpv.Cmd.Run(); err != nil {
		log.Println(err)
	}

	// Give us a 5-second period timeout.
	time.Sleep(time.Second)
	log.Println("mpv running")
	ipcc := remoteMpv.NewIPCClient(mpv.SocketPath) // Lowlevel client
	log.Println("ippc client created")
	c := remoteMpv.NewClient(ipcc) // Highlevel client, can also use RPCClient

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
		go play(c, input.Text)
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

func play(c *remoteMpv.Client, url string) {
	c.Loadfile(url, remoteMpv.LoadListModeReplace)
	c.SetPause(false)
}
