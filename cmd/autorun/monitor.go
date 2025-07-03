package main

import (
	"log"
	"time"

	driveutil "github.com/Merith-TK/utils/pkg/driveutil"
)

// runAutorunForDrive executes autorun functionality for a specific drive
func runAutorunForDrive(drive string) {
	log.Printf("[AUTORUN] Checking drive %s for autorun config", drive)
	startAutorun(drive)
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
