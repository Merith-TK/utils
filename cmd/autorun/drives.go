package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Merith-TK/utils/pkg/debug"
)

var driveStore = make(map[string]bool)

func detectDrives() {
	for _, drive := range "DEFGHIJKLMNOPQRSTUVWXYZ" {

		drivePath := string(drive) + ":\\"
		drives, err := filepath.Glob(drivePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, drive := range drives {
			if _, ok := driveStore[drive]; !ok {
				driveStore[drive] = true
				debug.Print("New drive detected:", drive, "\n")
				go startAutorun(drive)

			}
		}

		for drive := range driveStore {
			if !driveExists(drive) {
				delete(driveStore, drive)
				debug.Print("Drive removed: ", drive, "\n")
			}
		}
	}
}

func driveExists(drive string) bool {
	_, err := os.Stat(drive)
	return err == nil
}
