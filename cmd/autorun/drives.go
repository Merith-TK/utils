package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var driveStore = make(map[string]bool)

func detectDrives() {
	for _, drive := range "DEFGHIJKLMNOPQRSTUVWXYZ" {
		// for _, drive := range "E" {

		drivePath := string(drive) + ":\\"
		drives, err := filepath.Glob(drivePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, drive := range drives {
			if _, ok := driveStore[drive]; !ok {
				driveStore[drive] = true
				fmt.Printf("New drive detected: %s\n", drive)
				go startAutorun(drive)

			}
		}

		for drive := range driveStore {
			if !driveExists(drive) {
				delete(driveStore, drive)
				fmt.Printf("Drive removed: %s\n", drive)
			}
		}
	}
}

// func hasAutorunFile(drive string) bool {
// 	autorunFile := filepath.Join(drive, ".autorun.toml")
// 	_, err := os.Stat(autorunFile)
// 	return err == nil
// }

func driveExists(drive string) bool {
	_, err := os.Stat(drive)
	return err == nil
}
