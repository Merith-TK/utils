package main

import (
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Merith-TK/utils/pkg/config"
)

var (
	configWinMu    sync.Mutex
	openConfigWins = make(map[string]fyne.Window)
)

func showDriveConfigDialog(win fyne.Window, drive DriveInfo) {
	configWinMu.Lock()
	if existing, ok := openConfigWins[drive.Letter]; ok {
		configWinMu.Unlock()
		existing.RequestFocus()
		return
	}
	configWinMu.Unlock()
	title := "Create Autorun Config"
	if drive.HasConfig {
		title = "Edit Autorun Config"
	}
	configWin := fyne.CurrentApp().NewWindow(title)
	configWin.Resize(fyne.NewSize(400, 350))
	label := widget.NewLabel("Drive: " + drive.Letter + " (" + drive.Label + ")")
	autorunPath := filepath.Join(drive.Letter, ".autorun.toml")
	pathLabel := widget.NewLabel("Path: " + autorunPath)

	// Load config if present
	cfg := Config{Environment: map[string]string{}}
	if drive.HasConfig {
		config.LoadToml(&cfg, autorunPath)
	}

	autorunEntry := widget.NewEntry()
	autorunEntry.SetText(cfg.Autorun)
	workDirEntry := widget.NewEntry()
	workDirEntry.SetText(cfg.WorkDir)
	isolateCheck := widget.NewCheck("Isolate", nil)
	isolateCheck.SetChecked(cfg.Isolate)

	// Simple key-value editor for Environment
	envRows := []*widget.Entry{}
	keyRows := []*widget.Entry{}
	envBox := container.NewVBox()
	for k, v := range cfg.Environment {
		keyEntry := widget.NewEntry()
		keyEntry.SetText(k)
		valEntry := widget.NewEntry()
		valEntry.SetText(v)
		row := container.NewHBox(keyEntry, widget.NewLabel("="), valEntry)
		envBox.Add(row)
		keyRows = append(keyRows, keyEntry)
		envRows = append(envRows, valEntry)
	}
	addEnvBtn := widget.NewButton("Add Env", func() {
		keyEntry := widget.NewEntry()
		valEntry := widget.NewEntry()
		row := container.NewHBox(keyEntry, widget.NewLabel("="), valEntry)
		envBox.Add(row)
		keyRows = append(keyRows, keyEntry)
		envRows = append(envRows, valEntry)
		configWin.Content().Refresh()
	})

	saveBtn := widget.NewButton("Save", func() {
		cfg.Autorun = autorunEntry.Text
		cfg.WorkDir = workDirEntry.Text
		cfg.Isolate = isolateCheck.Checked
		cfg.Environment = map[string]string{}
		for i := range keyRows {
			k := keyRows[i].Text
			v := envRows[i].Text
			if k != "" {
				cfg.Environment[k] = v
			}
		}
		if err := config.SaveToml(autorunPath, cfg); err == nil {
			configWin.Close()
			configWinMu.Lock()
			delete(openConfigWins, drive.Letter)
			configWinMu.Unlock()
		} else {
			fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "Save Error", Content: err.Error()})
		}
	})
	cancelBtn := widget.NewButton("Cancel", func() {
		configWinMu.Lock()
		delete(openConfigWins, drive.Letter)
		configWinMu.Unlock()
		configWin.Close()
	})
	configWin.SetContent(container.NewVBox(
		label, pathLabel,
		widget.NewForm(
			widget.NewFormItem("Autorun", autorunEntry),
			widget.NewFormItem("WorkDir", workDirEntry),
			widget.NewFormItem("Isolate", isolateCheck),
		),
		widget.NewLabel("Environment:"),
		envBox, addEnvBtn,
		container.NewHBox(saveBtn, cancelBtn),
	))
	configWin.SetCloseIntercept(func() {
		configWinMu.Lock()
		delete(openConfigWins, drive.Letter)
		configWinMu.Unlock()
		configWin.Close()
	})
	configWinMu.Lock()
	openConfigWins[drive.Letter] = configWin
	configWinMu.Unlock()
	configWin.Show()
}
