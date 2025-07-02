//go:build windows

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/getlantern/systray"
)

var timeout int

func main() {
	flag.IntVar(&timeout, "timeout", 15, "Exit after N seconds (for testing)")
	flag.Parse()
	if timeout > 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			log.Printf("[TIMEOUT] Exiting after %d seconds (self-destruct)", timeout)
			os.Exit(0)
		}()
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("TrayTest")
	systray.SetTooltip("Fyne + Systray Example")
	// You can use your own icon here
	systray.SetIcon(nil)

	mOpen := systray.AddMenuItem("Open Window", "Show the Fyne window")
	mQuit := systray.AddMenuItem("Quit", "Exit the whole app")

	// Fyne app and window
	fyneApp := app.New()
	var win fyne.Window
	windowOpen := false

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				if !windowOpen {
					win = fyneApp.NewWindow("Fyne Test Window")
					win.SetContent(widget.NewLabel("Hello from Fyne!"))
					win.Resize(fyne.NewSize(300, 100))
					windowOpen = true
					win.SetCloseIntercept(func() {
						windowOpen = false
						win.Close()
					})
					go func() {
						win.ShowAndRun()
						windowOpen = false
					}()
				}
			case <-mQuit.ClickedCh:
				log.Println("[TRAY] Quit selected")
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	log.Println("[TRAY] onExit called, cleaning up...")
}
