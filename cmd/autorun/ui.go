package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	driveutil "github.com/Merith-TK/utils/pkg/driveutil"
)

// buildMainContent creates the main window content with drive listing
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

	return container.NewBorder(
		widget.NewLabelWithStyle("Detected Drives", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, // bottom
		nil, // left
		nil, // right
		container.NewStack(
			container.NewVScroll(table),
		),
	)
}
