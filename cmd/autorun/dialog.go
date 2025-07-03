package main

import (
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	configWin.Resize(fyne.NewSize(500, 600))
	configWin.SetFixedSize(true)
	
	// Load config if present
	cfg := Config{Environment: map[string]string{}}
	if drive.HasConfig {
		config.LoadToml(&cfg, filepath.Join(drive.Letter, ".autorun.toml"))
	}
	
	// Create the content
	content := createConfigDialogContent(drive, cfg, configWin)
	configWin.SetContent(content)
	
	// Handle window close
	configWin.SetCloseIntercept(func() {
		configWinMu.Lock()
		delete(openConfigWins, drive.Letter)
		configWinMu.Unlock()
		configWin.Close()
	})
	
	// Track the window
	configWinMu.Lock()
	openConfigWins[drive.Letter] = configWin
	configWinMu.Unlock()
	
	configWin.Show()
}

// createConfigDialogContent creates the modern content for the config dialog
func createConfigDialogContent(drive DriveInfo, cfg Config, configWin fyne.Window) fyne.CanvasObject {
	// Header section
	driveIcon := widget.NewIcon(theme.StorageIcon())
	driveTitle := widget.NewLabelWithStyle(drive.Letter, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	driveSubtitle := widget.NewLabelWithStyle(drive.Label, fyne.TextAlignLeading, fyne.TextStyle{})
	if drive.Label == "" {
		driveSubtitle.SetText("Unnamed Drive")
	}
	
	autorunPath := filepath.Join(drive.Letter, ".autorun.toml")
	pathLabel := widget.NewLabelWithStyle("Config Path: "+autorunPath, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	
	header := container.NewHBox(
		driveIcon,
		container.NewVBox(
			driveTitle,
			driveSubtitle,
			pathLabel,
		),
	)
	
	// Main configuration form
	autorunEntry := widget.NewEntry()
	autorunEntry.SetText(cfg.Autorun)
	autorunEntry.SetPlaceHolder("Path to executable (e.g., /setup.exe)")
	
	workDirEntry := widget.NewEntry()
	workDirEntry.SetText(cfg.WorkDir)
	workDirEntry.SetPlaceHolder("Working directory (optional)")
	
	isolateCheck := widget.NewCheck("Enable Isolation", nil)
	isolateCheck.SetChecked(cfg.Isolate)
	
	// Create help text for isolation
	isolateHelp := widget.NewRichTextFromMarkdown(`
**Isolation Mode**: When enabled, the application runs in a sandboxed environment with limited access to your system. This provides better security but may prevent some applications from working properly.`)
	isolateHelp.Wrapping = fyne.TextWrapWord
	
	// Main form
	form := container.NewVBox(
		widget.NewCard("Basic Configuration", "", container.NewVBox(
			widget.NewLabelWithStyle("Autorun Command", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			autorunEntry,
			widget.NewLabelWithStyle("Working Directory", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			workDirEntry,
			isolateCheck,
			isolateHelp,
		)),
	)
	
	// Environment variables section
	envCard, envRows := createEnvironmentCard(cfg.Environment)
	form.Add(envCard)
	
	// Action buttons
	saveBtn := widget.NewButton("Save Configuration", func() {
		saveConfig(drive, autorunEntry, workDirEntry, isolateCheck, envRows, configWin)
	})
	saveBtn.Importance = widget.HighImportance
	
	cancelBtn := widget.NewButton("Cancel", func() {
		configWinMu.Lock()
		delete(openConfigWins, drive.Letter)
		configWinMu.Unlock()
		configWin.Close()
	})
	
	buttonContainer := container.NewHBox(
		cancelBtn,
		saveBtn,
	)
	
	// Scroll container for the form
	scrollContainer := container.NewVScroll(form)
	
	// Main layout
	return container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()),
		buttonContainer,
		nil,
		nil,
		scrollContainer,
	)
}

// createEnvironmentCard creates a card for environment variables
func createEnvironmentCard(environment map[string]string) (*widget.Card, *[]*envRow) {
	envRows := []*envRow{}
	envContainer := container.NewVBox()
	
	// Add existing environment variables
	for k, v := range environment {
		row := createEnvironmentRow(k, v, &envRows, envContainer)
		envContainer.Add(row)
	}
	
	// Add environment variable button
	addEnvBtn := widget.NewButton("Add Environment Variable", func() {
		row := createEnvironmentRow("", "", &envRows, envContainer)
		envContainer.Add(row)
		envContainer.Refresh()
	})
	addEnvBtn.SetIcon(theme.ContentAddIcon())
	
	cardContent := container.NewVBox(
		envContainer,
		addEnvBtn,
	)
	
	return widget.NewCard("Environment Variables", "Custom environment variables for the application", cardContent), &envRows
}

// envRow represents a row of environment variable inputs
type envRow struct {
	keyEntry   *widget.Entry
	valueEntry *widget.Entry
	container  *fyne.Container
}

// createEnvironmentRow creates a row for editing environment variables
func createEnvironmentRow(key, value string, envRows *[]*envRow, envContainer *fyne.Container) *fyne.Container {
	keyEntry := widget.NewEntry()
	keyEntry.SetText(key)
	keyEntry.SetPlaceHolder("Variable name")
	
	valueEntry := widget.NewEntry()
	valueEntry.SetText(value)
	valueEntry.SetPlaceHolder("Variable value")
	
	deleteBtn := widget.NewButton("", func() {
		// Remove this row
		for i, row := range *envRows {
			if row.keyEntry == keyEntry {
				envContainer.Remove(row.container)
				*envRows = append((*envRows)[:i], (*envRows)[i+1:]...)
				envContainer.Refresh()
				break
			}
		}
	})
	deleteBtn.SetIcon(theme.DeleteIcon())
	deleteBtn.Importance = widget.DangerImportance
	
	rowContainer := container.NewBorder(
		nil, nil, nil, deleteBtn,
		container.NewGridWithColumns(2, keyEntry, valueEntry),
	)
	
	row := &envRow{
		keyEntry:   keyEntry,
		valueEntry: valueEntry,
		container:  rowContainer,
	}
	*envRows = append(*envRows, row)
	
	return rowContainer
}

// saveConfig saves the configuration
func saveConfig(drive DriveInfo, autorunEntry, workDirEntry *widget.Entry, isolateCheck *widget.Check, envRows *[]*envRow, configWin fyne.Window) {
	cfg := Config{
		Autorun:     autorunEntry.Text,
		WorkDir:     workDirEntry.Text,
		Isolate:     isolateCheck.Checked,
		Environment: make(map[string]string),
	}
	
	// Extract environment variables from the rows
	for _, row := range *envRows {
		key := row.keyEntry.Text
		value := row.valueEntry.Text
		if key != "" {
			cfg.Environment[key] = value
		}
	}
	
	autorunPath := filepath.Join(drive.Letter, ".autorun.toml")
	if err := config.SaveToml(autorunPath, cfg); err == nil {
		configWin.Close()
		configWinMu.Lock()
		delete(openConfigWins, drive.Letter)
		configWinMu.Unlock()
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Configuration Saved",
			Content: "Autorun configuration saved successfully",
		})
	} else {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Save Error",
			Content: err.Error(),
		})
	}
}
