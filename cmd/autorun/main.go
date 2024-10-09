package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Merith-TK/utils/debug"
)

func main() {
	// check if "install" argument is provided
	if len(os.Args) > 1 && os.Args[1] == "install" {
		debug.Print("Installing autorun service")
		copyToStartupFolder()
		return
	}

	// check if .autorun.toml file exists next to the executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(filepath.Dir(exePath), ".autorun.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		debug.Print("Starting as autorun service")
		for {
			detectDrives()
			time.Sleep(5 * time.Second)
		}
	} else {
		debug.Print("Starting as autorun program")
		startAutorun(filepath.Dir(exePath))
	}
	line := ""
	fmt.Printf("Invalid glob pattern '%s' (skipped): %v\n", line, err)
	fmt.Printf("util.RemoveAll: %v\n", err)

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
