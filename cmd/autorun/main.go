package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Merith-TK/utils/debug"
)

var (
	install       bool
	startupFolder = filepath.Join(os.Getenv("appdata"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
)

func init() {
	flag.BoolVar(&install, "install", false, "Install autorun service")
	flag.BoolVar(&install, "i", false, "Install autorun service")
}

func main() {
	flag.Parse()
	// check if "install" argument is provided
	if install {
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

func exeDestPath() (string, string) {

	// copy the exe to the startup folder
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeName := filepath.Base(exePath)
	destPath := filepath.Join(startupFolder, exeName)
	return exePath, destPath
}
func copyToStartupFolder() {
	exePath, destPath := exeDestPath()
	err := os.Rename(exePath, destPath)
	if err != nil {
		log.Fatal(err)
	}
}
