// Package main implements a test utility for the driveutil package functionality,\r\n// demonstrating drive detection and enumeration capabilities.\r\n//\r\n// The testdriveutil utility provides a simple command-line interface to test\r\n// and demonstrate the driveutil package's drive detection functionality.\r\n// It lists all available drives with their metadata including labels,\r\n// serial numbers, and drive types.\r\n//\r\n// Features:\r\n//   - Drive enumeration using driveutil.ListDrives()\r\n//   - Display of drive metadata (letter, label, serial, type)\r\n//   - Simple output formatting for easy reading\r\n//   - Detection of no-drive scenarios\r\n//\r\n// Usage:\r\n//   testdriveutil\r\n//\r\n// Output Format:\r\n//   Drive_Letter    Label: Volume_Label    Serial: XXXXXXXX    Type: N\r\n//\r\n// Example Output:\r\n//   Detected drives:\r\n//   C:\\    Label: Windows    Serial: 12345678    Type: 3\r\n//   D:\\    Label: Data       Serial: 87654321    Type: 3\r\n//   E:\\    Label: USB Drive  Serial: ABCDEF00    Type: 2\r\n//\r\n// This utility is primarily used for testing and debugging the driveutil\r\n// package functionality on different Windows systems.\r\npackage main

import (
	"fmt"

	driveutil "github.com/Merith-TK/utils/pkg/driveutil"
)

func main() {
	drives := driveutil.ListDrives()
	if len(drives) == 0 {
		fmt.Println("No drives detected.")
		return
	}
	fmt.Println("Detected drives:")
	for _, d := range drives {
		fmt.Printf("%s\tLabel: %s\tSerial: %08X\tType: %d\n", d.Letter, d.Label, d.Serial, d.Type)
	}
}
