package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
)

// onReady sets up the system tray menu and handles menu events
func onReady(showWinCh chan struct{}, quitCh chan struct{}) {
	systray.SetTitle("Autorun")
	systray.SetTooltip("Autorun Manager")
	// systray.SetIcon(nil) // Add icon if desired

	mOpen := systray.AddMenuItem("Open GUI", "Show the Fyne window")
	mQuit := systray.AddMenuItem("Quit", "Exit the whole app")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				log.Println("[TRAY] Open GUI clicked")
				showWinCh <- struct{}{}
			case <-mQuit.ClickedCh:
				log.Println("[TRAY] Quit clicked")
				quitCh <- struct{}{}
				os.Exit(0)
			}
		}
	}()
}

// onExit handles cleanup when the system tray exits
func onExit() {
	log.Println("[TRAY] onExit called, cleaning up...")
}
