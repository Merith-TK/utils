package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	driveutil "github.com/Merith-TK/utils/pkg/driveutil"
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

type DriveInfo struct {
	Letter    string
	Label     string
	HasConfig bool
}

func buildMainContent(win fyne.Window, configDialogCh chan<- DriveInfo) fyne.CanvasObject {
	drives := []DriveInfo{}
	for _, d := range driveutil.ListDrives() {
		// Check for `.autorun.toml` on the root of the drive
		autorunPath := filepath.Join(d.Letter, ".autorun.toml")
		hasConfig := false
		if _, err := os.Stat(autorunPath); err == nil {
			hasConfig = true
		}
		drives = append(drives, DriveInfo{
			Letter:    d.Letter,
			Label:     d.Label,
			HasConfig: hasConfig,
		})
	}
	columns := []string{"Drive", "Label", "Has Config"}
	table := widget.NewTable(
		func() (int, int) { return len(drives) + 1, len(columns) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.SetText(columns[id.Col])
			} else {
				drive := drives[id.Row-1]
				switch id.Col {
				case 0:
					label.SetText(drive.Letter)
				case 1:
					label.SetText(drive.Label)
				case 2:
					if drive.HasConfig {
						label.SetText("Yes")
					} else {
						label.SetText("No")
					}
				}
				label.TextStyle = fyne.TextStyle{} // not bold for data
			}
		},
	)
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			configDialogCh <- drives[id.Row-1]
		}
	}
	table.SetColumnWidth(0, 60)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 80)

	// Create a border layout with the title at the top and table taking the rest
	return container.NewBorder(
		widget.NewLabelWithStyle("Detected Drives", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, // bottom
		nil, // left
		nil, // right
		container.NewStack( // This will make the table expand to fill available space
			container.NewVScroll(table),
		),
	)
}

type uiAction func()

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

func onExit() {
	log.Println("[TRAY] onExit called, cleaning up...")
}

func exeDestPath() (string, string) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("[INSTALL ERROR] Failed to get executable path:", err)
	}
	exeName := filepath.Base(exePath)
	destPath := filepath.Join(startupFolder, exeName)
	return exePath, destPath
}
func copyToStartupFolder() {
	exePath, destPath := exeDestPath()
	err := os.Rename(exePath, destPath)
	if err != nil {
		log.Fatal("[INSTALL ERROR] Failed to move exe to startup folder:", err)
	}
}

func runAutorunForDrive(drive string) {
	log.Printf("[AUTORUN] Checking drive %s for autorun config", drive)
	startAutorun(drive)
}

func startDriveMonitor(uiRefreshCh chan<- struct{}) {
	go func() {
		log.Println("[MONITOR] Starting drive monitor")

		// Check existing drives on startup
		log.Println("[MONITOR] Checking existing drives...")
		for _, drive := range driveutil.ListDrives() {
			log.Printf("[MONITOR] Existing drive found: %s", drive.Letter)
			runAutorunForDrive(drive.Letter)
		}

		driveStore := driveutil.DriveStore{}
		driveStore.MonitorDrives(func(drive string, serial uint32) {
			log.Printf("[MONITOR] New drive detected: %s (serial: %d)", drive, serial)
			runAutorunForDrive(drive)
			if uiRefreshCh != nil {
				uiRefreshCh <- struct{}{}
			}
		}, 5*time.Second)
	}()
}
