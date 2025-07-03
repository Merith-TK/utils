package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/getlantern/systray"
)

var (
	install       bool
	timeout       int // seconds, 0 means no timeout
	startupFolder = filepath.Join(os.Getenv("appdata"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
)

func init() {
	flag.BoolVar(&install, "install", false, "Install autorun service")
	flag.BoolVar(&install, "i", false, "Install autorun service")
	flag.IntVar(&timeout, "timeout", 0, "Exit after N seconds (for testing)")
}

func main() {
	flag.Parse()
	log.Printf("[MAIN] Flags parsed. install=%v, timeout=%v", install, timeout)
	if timeout > 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			log.Printf("[TIMEOUT] Exiting after %d seconds (self-destruct)", timeout)
			os.Exit(0)
		}()
	}
	if install {
		copyToStartupFolder()
		return
	}

	showWinCh := make(chan struct{}, 1)
	quitCh := make(chan struct{}, 1)
	configDialogCh := make(chan DriveInfo, 1)
	uiActionCh := make(chan uiAction, 10)
	uiRefreshCh := make(chan struct{}, 1)

	go func() {
		systray.Run(func() { onReady(showWinCh, quitCh) }, onExit)
	}()

	fyneApp := app.New()
	win := fyneApp.NewWindow("Autorun Drive Manager")
	win.Resize(fyne.NewSize(400, 300))
	win.SetFixedSize(true)
	win.SetCloseIntercept(func() {
		log.Println("[FYNE] Window close intercepted, hiding window")
		win.Hide()
	})

	refreshContent := func() {
		win.SetContent(buildMainContent(win, configDialogCh))
	}

	refreshContent() // initial content

	// Tray/UI event goroutine: send UI actions to uiActionCh
	go func() {
		for {
			select {
			case <-showWinCh:
				log.Println("[FYNE] Show window requested from tray")
				uiActionCh <- func() {
					refreshContent()
					win.Show()
					win.RequestFocus()
				}
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "UI", Content: "Show window"})
			case <-quitCh:
				log.Println("[FYNE] Quit requested from tray")
				uiActionCh <- func() { fyneApp.Quit() }
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "UI", Content: "Quit"})
				return
			case drive := <-configDialogCh:
				uiActionCh <- func() { showDriveConfigDialog(win, drive) }
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "UI", Content: "Config dialog"})
			}
		}
	}()

	// Timer-based polling for UI actions on the main thread
	go func() {
		for {
			select {
			case action := <-uiActionCh:
				action()
			case <-uiRefreshCh:
				uiActionCh <- refreshContent
			case <-time.After(50 * time.Millisecond):
				// allow UI to update
			}
		}
	}()

	startDriveMonitor(uiRefreshCh)

	fyneApp.Run()
}


