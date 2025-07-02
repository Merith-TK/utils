package main

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
