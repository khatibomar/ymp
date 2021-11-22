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

	// Give us a 5-second period timeout.
	time.Sleep(time.Second)
	log.Println("mpv running")
	ipcc := remoteMpv.NewIPCClient(mpv.SocketPath) // Lowlevel client
	log.Println("ippc client created")
	c := remoteMpv.NewClient(ipcc) // Highlevel client, can also use RPCClient

	// components declaration
	var input *widget.Entry
	var slider *widget.Slider
	var playBtn *widget.Button
	var pauseBtn *widget.Button
	var logBox *widget.TextGrid
	var logLabel *widget.Label
	var logText string

	paused := false

	input = widget.NewEntry()
	input.SetPlaceHolder("Enter youtube URL...")

	slider = widget.NewSlider(0, 100)
	slider.Value = 100
	slider.OnChanged = func(f float64) {
		c.SetProperty("volume", f)
	}

	linkValidator := validation.NewRegexp(ytExp, "Invalid youtube link")

	playBtn = widget.NewButton("Play", func() {
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
		c.Loadfile(input.Text, remoteMpv.LoadListModeReplace)
		c.SetPause(false)
		pauseBtn.Enable()
		logText = ""
		logBox.SetText(logText)
	})

	pauseBtn = widget.NewButton("Pause", func() {
		log.Println("pause called")
		paused = !paused
		if paused {
			pauseBtn.Text = "Resume"
		} else {
			pauseBtn.Text = "Pause"
		}
		c.SetPause(paused)
		pauseBtn.Refresh()
	})
	pauseBtn.Disable()

	logBox = widget.NewTextGrid()

	go func() {
		for {
			line, _, _ := mpv.OutBuff.ReadLine()
			if string(line) == "\n" || string(line) == "" {
				continue
			}
			logText += string(line) + "\n"
			logBox.SetText(logText)
		}
	}()

	logLabel = widget.NewLabel("Logs:")

	content := container.New(
		layout.NewVBoxLayout(),
		input,
		playBtn,
		pauseBtn,
		slider,
		logLabel,
		logBox,
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
