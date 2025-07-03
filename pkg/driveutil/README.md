# driveutil

This package provides drive enumeration, metadata extraction, and utility functions for Windows drives.

## Types

### DriveStore

```
type DriveStore map[string]bool
```
Tracks detected drives by unique ID.

### DriveInfo

```
type DriveInfo struct {
    Letter string
    Label  string
    Serial uint32
    Type   uint32 // DRIVE_FIXED, DRIVE_REMOVABLE, etc.
}
```

## Functions

- `(store DriveStore) DetectDrives(onNewDrive func(drive string, serial uint32))` - Detects new drives and calls callback
- `(store DriveStore) MonitorDrives(onNewDrive func(drive string, serial uint32), interval time.Duration)` - Continuously monitors for new drives
- `GetVolumeSerialNumber(root string) (uint32, error)` - Gets volume serial number for drive
- `DriveExists(drive string) bool` - Checks if drive path exists
- `ListDrives() []DriveInfo` - Returns slice of all available drives

## Example

```go
import "github.com/Merith-TK/utils/pkg/driveutil"

drives := driveutil.ListDrives()
for _, drive := range drives {
    fmt.Printf("Drive: %s, Label: %s, Serial: %08X\n", drive.Letter, drive.Label, drive.Serial)
}

store := driveutil.DriveStore{}
store.MonitorDrives(func(drive string, serial uint32) {
    fmt.Printf("New drive detected: %s (Serial: %08X)\n", drive, serial)
}, 5*time.Second)
``` 