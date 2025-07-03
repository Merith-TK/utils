package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Merith-TK/utils/pkg/driveutil"
)

var (
	securityDialogMu    sync.Mutex
	openSecurityDialogs = make(map[string]fyne.Window)
)

// SecurityDialogResult represents the result of a security dialog
type SecurityDialogResult struct {
	Decision SecurityDecision
	Remember bool
}

// showSecurityDialog shows a security dialog for an unknown or changed config
func showSecurityDialog(metadata *ConfigMetadata, drivePath string) (*SecurityDialogResult, error) {
	// Get drive serial for dialog tracking
	driveSerial, err := driveutil.GetVolumeSerialNumber(drivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get drive serial: %v", err)
	}
	
	// Prevent multiple dialogs for the same drive
	securityDialogMu.Lock()
	dialogKey := fmt.Sprintf("%08X", driveSerial)
	if existingDialog, exists := openSecurityDialogs[dialogKey]; exists {
		securityDialogMu.Unlock()
		existingDialog.RequestFocus()
		return nil, fmt.Errorf("dialog already open for this drive")
	}
	securityDialogMu.Unlock()

	// Create the dialog window
	app := fyne.CurrentApp()
	dialog := app.NewWindow("Security Warning - Autorun Config Detected")
	dialog.Resize(fyne.NewSize(550, 450))
	dialog.SetFixedSize(true)

	// Track the dialog
	securityDialogMu.Lock()
	openSecurityDialogs[dialogKey] = dialog
	securityDialogMu.Unlock()

	// Result channel
	resultCh := make(chan *SecurityDialogResult, 1)

	// Cleanup function
	cleanup := func() {
		securityDialogMu.Lock()
		delete(openSecurityDialogs, dialogKey)
		securityDialogMu.Unlock()
		dialog.Close()
	}

	// Create the content
	content := createSecurityDialogContent(metadata, drivePath, resultCh, cleanup)
	dialog.SetContent(content)

	// Handle window close
	dialog.SetCloseIntercept(func() {
		cleanup()
		resultCh <- &SecurityDialogResult{Decision: SecurityDecisionDeny, Remember: false}
	})

	// Show the dialog
	dialog.Show()
	dialog.RequestFocus()

	// Wait for result
	result := <-resultCh
	return result, nil
}

// createSecurityDialogContent creates the content for the security dialog
func createSecurityDialogContent(metadata *ConfigMetadata, drivePath string, resultCh chan<- *SecurityDialogResult, cleanup func()) fyne.CanvasObject {
	// Warning header
	warningLabel := widget.NewLabelWithStyle("⚠️ SECURITY WARNING", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	warningLabel.Wrapping = fyne.TextWrapWord

	// Drive information
	driveInfoLabel := widget.NewLabelWithStyle(fmt.Sprintf("Drive: %s", drivePath), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Hash information (on one line)
	hashInfoLabel := widget.NewLabelWithStyle(fmt.Sprintf("Config Hash (MD5): %s", metadata.MD5Hash), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	hashInfoLabel.TextStyle = fyne.TextStyle{Monospace: true}

	// Config details
	configLabel := widget.NewLabelWithStyle("Configuration Details:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	configDetails := createConfigDetailsWidget(metadata)

	// Environment variables
	envLabel := widget.NewLabelWithStyle("Environment Variables:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	envDetails := createEnvironmentDetailsWidget(metadata)

	// Remember decision checkbox
	rememberCheck := widget.NewCheck("Remember my decision", nil)
	rememberCheck.SetChecked(true)

	// Buttons
	allowBtn := widget.NewButton("Allow", func() {
		decision := SecurityDecisionAllow
		if !rememberCheck.Checked {
			decision = SecurityDecisionAllowOnce
		}
		cleanup()
		resultCh <- &SecurityDialogResult{Decision: decision, Remember: rememberCheck.Checked}
	})
	allowBtn.Importance = widget.SuccessImportance

	denyBtn := widget.NewButton("Deny", func() {
		decision := SecurityDecisionDeny
		if !rememberCheck.Checked {
			decision = SecurityDecisionDenyOnce
		}
		cleanup()
		resultCh <- &SecurityDialogResult{Decision: decision, Remember: rememberCheck.Checked}
	})
	denyBtn.Importance = widget.DangerImportance

	// Layout
	// Header section (non-scrollable) - only warning, drive, and hash
	headerSection := container.NewVBox(
		warningLabel,
		widget.NewSeparator(),
		container.NewHBox(driveInfoLabel),
		container.NewHBox(hashInfoLabel),
	)

	// Scrollable content - only config details and environment
	scrollableContent := container.NewVBox(
		configLabel,
		configDetails,
		envLabel,
		envDetails,
	)

	// Footer section (non-scrollable) - checkbox and buttons
	footerSection := container.NewVBox(
		widget.NewSeparator(),
		rememberCheck,
		container.NewHBox(
			denyBtn,
			widget.NewSeparator(),
			allowBtn,
		),
	)

	return container.NewBorder(
		headerSection,
		footerSection,
		nil,
		nil,
		container.NewScroll(scrollableContent),
	)
}

// createConfigDetailsWidget creates a widget showing config details
func createConfigDetailsWidget(metadata *ConfigMetadata) fyne.CanvasObject {
	cfg := metadata.Config

	details := []string{}

	if cfg.Autorun != "" {
		details = append(details, fmt.Sprintf("• Command: %s", cfg.Autorun))
	}

	if cfg.WorkDir != "" {
		details = append(details, fmt.Sprintf("• Working Directory: %s", cfg.WorkDir))
	}

	if cfg.Isolate {
		details = append(details, "• Isolation: Enabled (sandboxed environment)")
	} else {
		details = append(details, "• Isolation: Disabled (full system access)")
	}

	if len(details) == 0 {
		details = append(details, "• No configuration details available")
	}

	text := strings.Join(details, "\n")
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		nil, nil, widget.NewLabel("  "), nil,
		label,
	)
}

// createEnvironmentDetailsWidget creates a widget showing environment variables
func createEnvironmentDetailsWidget(metadata *ConfigMetadata) fyne.CanvasObject {
	env := metadata.Environment

	if len(env) == 0 {
		label := widget.NewLabel("• No custom environment variables")
		return container.NewBorder(
			nil, nil, widget.NewLabel("  "), nil,
			label,
		)
	}

	// Sort environment variables by key
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	details := []string{}
	for _, key := range keys {
		value := env[key]
		if len(value) > 50 {
			value = value[:47] + "..."
		}
		details = append(details, fmt.Sprintf("• %s = %s", key, value))
	}

	text := strings.Join(details, "\n")
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle = fyne.TextStyle{Monospace: true}

	return container.NewBorder(
		nil, nil, widget.NewLabel("  "), nil,
		label,
	)
}
