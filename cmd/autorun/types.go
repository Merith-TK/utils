package main

// DriveInfo represents information about a drive
type DriveInfo struct {
	Letter    string
	Label     string
	HasConfig bool
}

// uiAction represents a UI action function
type uiAction func()
