package main

import (
	"log"
	"time"

	driveutil "github.com/Merith-TK/utils/pkg/driveutil"
)

// runAutorunForDrive executes autorun functionality for a specific drive
func runAutorunForDrive(drive string) {
	log.Printf("[AUTORUN] Checking drive %s for autorun config", drive)

	// Check security decision
	decision, metadata, err := securityManager.CheckConfig(drive)
	if err != nil {
		log.Printf("[SECURITY] Error checking config security: %v", err)
		return
	}

	switch decision {
	case SecurityDecisionAllow, SecurityDecisionAllowOnce:
		log.Printf("[SECURITY] Config approved for drive %s", drive)
		startAutorun(drive)
	case SecurityDecisionDeny, SecurityDecisionDenyOnce:
		log.Printf("[SECURITY] Config denied for drive %s", drive)
		return
	case SecurityDecisionUnknown:
		log.Printf("[SECURITY] Unknown config detected for drive %s, showing security dialog", drive)

		// Show security dialog
		result, err := showSecurityDialog(metadata, drive)
		if err != nil {
			log.Printf("[SECURITY] Error showing security dialog: %v", err)
			return
		}

		// Save the decision
		err = securityManager.SaveDecision(metadata, result.Decision, drive)
		if err != nil {
			log.Printf("[SECURITY] Error saving decision: %v", err)
			return
		}

		// Act on the decision
		switch result.Decision {
		case SecurityDecisionAllow, SecurityDecisionAllowOnce:
			log.Printf("[SECURITY] User approved config for drive %s", drive)
			startAutorun(drive)
		case SecurityDecisionDeny, SecurityDecisionDenyOnce:
			log.Printf("[SECURITY] User denied config for drive %s", drive)
			return
		}
	}
}

// startDriveMonitor starts monitoring for drive changes and executes autorun
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
