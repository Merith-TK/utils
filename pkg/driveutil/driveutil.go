// Package driveutil provides drive enumeration, metadata extraction, and utility functions for Windows drives.
package driveutil

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

// DriveStore keeps track of currently detected drives by their unique ID.
type DriveStore map[string]bool

// DriveInfo represents information about a drive.
type DriveInfo struct {
	Letter string
	Label  string
	Serial uint32
	Type   uint32 // windows.DRIVE_FIXED, DRIVE_REMOVABLE, etc.
}

// DetectDrives enumerates all logical drives, checks their type, and gets their serial number.
// Calls the provided callback for each new drive detected.
func (store DriveStore) DetectDrives(onNewDrive func(drive string, serial uint32)) {
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		log.Printf("Failed to get logical drives: %v", err)
		return
	}
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		drive := fmt.Sprintf("%c:\\", 'A'+i)
		driveType := windows.GetDriveType(syscall.StringToUTF16Ptr(drive))
		if driveType != windows.DRIVE_REMOVABLE && driveType != windows.DRIVE_FIXED {
			continue
		}
		serial, err := GetVolumeSerialNumber(drive)
		if err != nil {
			continue
		}
		uniqueID := fmt.Sprintf("%s-%08X", drive, serial)
		if _, ok := store[uniqueID]; !ok {
			store[uniqueID] = true
			onNewDrive(drive, serial)
		}
	}
	// Remove drives that are no longer present
	for uniqueID := range store {
		drive := uniqueID[:3]
		if !DriveExists(drive) {
			delete(store, uniqueID)
		}
	}
}

// GetVolumeSerialNumber returns the serial number for a given drive root.
func GetVolumeSerialNumber(root string) (uint32, error) {
	var (
		volumeName      [windows.MAX_PATH + 1]uint16
		fsName          [windows.MAX_PATH + 1]uint16
		serialNumber    uint32
		maxComponentLen uint32
		fileSystemFlags uint32
	)
	rootPtr, _ := syscall.UTF16PtrFromString(root)
	ret := windows.GetVolumeInformation(
		rootPtr,
		&volumeName[0],
		uint32(len(volumeName)),
		&serialNumber,
		&maxComponentLen,
		&fileSystemFlags,
		&fsName[0],
		uint32(len(fsName)),
	)
	if ret != nil {
		return 0, ret
	}
	return serialNumber, nil
}

// DriveExists checks if a drive path exists.
func DriveExists(drive string) bool {
	_, err := syscall.GetFileAttributes(syscall.StringToUTF16Ptr(drive))
	return err == nil
}

// MonitorDrives calls DetectDrives in a loop with a sleep interval.
func (store DriveStore) MonitorDrives(onNewDrive func(drive string, serial uint32), interval time.Duration) {
	for {
		store.DetectDrives(onNewDrive)
		time.Sleep(interval)
	}
}

// ListDrives returns a slice of DriveInfo for all present fixed/removable drives.
func ListDrives() []DriveInfo {
	var drives []DriveInfo
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		log.Printf("Failed to get logical drives: %v", err)
		return drives
	}
	log.Printf("Logical drive mask: %08b", mask)
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		drive := fmt.Sprintf("%c:\\", 'A'+i)
		log.Printf("Checking drive: %s", drive)
		ptr, _ := syscall.UTF16PtrFromString(drive)
		driveType := windows.GetDriveType(ptr)
		log.Printf("Drive %s type: %d", drive, driveType)
		if driveType != windows.DRIVE_REMOVABLE && driveType != windows.DRIVE_FIXED {
			log.Printf("Drive %s skipped (not removable/fixed)", drive)
			continue
		}
		serial, err := GetVolumeSerialNumber(drive)
		if err != nil {
			log.Printf("Drive %s serial error: %v", drive, err)
			continue
		}
		var volumeName [windows.MAX_PATH + 1]uint16
		var fsName [windows.MAX_PATH + 1]uint16
		err = windows.GetVolumeInformation(
			ptr,
			&volumeName[0],
			uint32(len(volumeName)),
			new(uint32), new(uint32), new(uint32),
			&fsName[0], uint32(len(fsName)),
		)
		if err != nil {
			log.Printf("Drive %s label error: %v", drive, err)
			continue
		}
		label := syscall.UTF16ToString(volumeName[:])
		log.Printf("Drive %s label: %s, serial: %08X", drive, label, serial)
		drives = append(drives, DriveInfo{
			Letter: drive,
			Label:  label,
			Serial: serial,
			Type:   driveType,
		})
	}
	log.Printf("Drives found: %d", len(drives))
	return drives
}
