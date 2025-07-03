package main

import (
	"log"
	"os"
	"path/filepath"
)

// exeDestPath returns the current executable path and the destination path in startup folder
func exeDestPath() (string, string) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("[INSTALL ERROR] Failed to get executable path:", err)
	}
	exeName := filepath.Base(exePath)
	destPath := filepath.Join(startupFolder, exeName)
	return exePath, destPath
}

// copyToStartupFolder copies the current executable to the Windows startup folder
func copyToStartupFolder() {
	exePath, destPath := exeDestPath()
	err := os.Rename(exePath, destPath)
	if err != nil {
		log.Fatal("[INSTALL ERROR] Failed to move exe to startup folder:", err)
	}
}
