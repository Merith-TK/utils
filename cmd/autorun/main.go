package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// check if "install" argument is provided
	if len(os.Args) > 1 && os.Args[1] == "install" {
		copyToStartupFolder()
		return
	}

	log.Println("Starting autorun service")
	for {
		detectDrives()
		time.Sleep(5 * time.Second)

	}
}

func copyToStartupFolder() {
	// get the path to the startup folder
	startupFolder, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	startupFolder = filepath.Join(startupFolder, "Microsoft", "Windows", "Start Menu", "Programs", "Startup")

	// copy the exe to the startup folder
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeName := filepath.Base(exePath)
	destPath := filepath.Join(startupFolder, exeName)
	err = os.Rename(exePath, destPath)
	if err != nil {
		log.Fatal(err)
	}
}
