package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	
	// Create header with icon and title
	headerIcon := widget.NewIcon(theme.StorageIcon())
	headerTitle := widget.NewLabelWithStyle("Autorun Drive Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	headerSubtitle := widget.NewLabelWithStyle("Detected Storage Drives", fyne.TextAlignCenter, fyne.TextStyle{})
	
	header := container.NewVBox(
		container.NewHBox(
			headerIcon,
			container.NewVBox(headerTitle, headerSubtitle),
		),
		widget.NewSeparator(),
	)
	
	// Create drive cards instead of table
	driveCards := container.NewVBox()
	
	if len(drives) == 0 {
		emptyIcon := widget.NewIcon(theme.InfoIcon())
		emptyLabel := widget.NewLabelWithStyle("No drives detected", fyne.TextAlignCenter, fyne.TextStyle{})
		emptyCard := container.NewBorder(
			nil, nil, emptyIcon, nil,
			emptyLabel,
		)
		driveCards.Add(emptyCard)
	} else {
		for i, drive := range drives {
			card := createDriveCard(drive, configDialogCh)
			driveCards.Add(card)
			
			// Add spacing between cards (except for the last one)
			if i < len(drives)-1 {
				driveCards.Add(widget.NewSeparator())
			}
		}
	}
	
	// Create footer with information
	footer := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Click on a drive to configure autorun settings", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
	)
	
	// Scroll container for the cards
	scrollContainer := container.NewVScroll(driveCards)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))
	
	return container.NewBorder(
		header,
		footer,
		nil,
		nil,
		scrollContainer,
	)
}

// createDriveCard creates a modern card for each drive
func createDriveCard(drive DriveInfo, configDialogCh chan<- DriveInfo) fyne.CanvasObject {
	// Drive type icon
	var driveIcon fyne.Resource
	switch drive.Letter {
	case "C:\\":
		driveIcon = theme.ComputerIcon()
	default:
		driveIcon = theme.StorageIcon()
	}
	
	icon := widget.NewIcon(driveIcon)
	icon.Resize(fyne.NewSize(32, 32))
	
	// Drive letter (large and bold)
	driveLabel := widget.NewLabelWithStyle(drive.Letter, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Drive name/label
	var displayLabel string
	if drive.Label == "" {
		displayLabel = "Unnamed Drive"
	} else {
		displayLabel = drive.Label
	}
	nameLabel := widget.NewLabelWithStyle(displayLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	// Config status with icon
	var configIcon fyne.Resource
	var configText string
	var configColor fyne.TextStyle
	
	if drive.HasConfig {
		configIcon = theme.ConfirmIcon()
		configText = "Autorun Config Found"
		configColor = fyne.TextStyle{}
	} else {
		configIcon = theme.InfoIcon()
		configText = "No Autorun Config"
		configColor = fyne.TextStyle{Italic: true}
	}
	
	configStatusIcon := widget.NewIcon(configIcon)
	configStatusLabel := widget.NewLabelWithStyle(configText, fyne.TextAlignLeading, configColor)
	
	// Action button
	actionButton := widget.NewButton("Configure", func() {
		configDialogCh <- drive
	})
	
	if drive.HasConfig {
		actionButton.SetText("Edit Config")
		actionButton.Importance = widget.MediumImportance
	} else {
		actionButton.SetText("Create Config")
		actionButton.Importance = widget.LowImportance
	}
	
	// Left side: icon and drive letter in a centered column
	leftSide := container.NewVBox(
		container.NewCenter(icon),
		container.NewCenter(driveLabel),
	)
	leftSide.Resize(fyne.NewSize(80, 60))
	
	// Middle section: drive info
	middleSection := container.NewVBox(
		nameLabel,
		container.NewHBox(
			configStatusIcon,
			configStatusLabel,
		),
	)
	
	// Right side: action button
	rightSide := container.NewVBox(
		actionButton,
	)
	
	// Main card content
	cardContent := container.NewBorder(
		nil, // top
		nil, // bottom
		leftSide, // left
		rightSide, // right
		middleSection, // center
	)
	
	// Add padding and border styling
	cardWithPadding := container.NewPadded(cardContent)
	
	return cardWithPadding
}
