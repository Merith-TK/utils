package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Merith-TK/utils/pkg/archive"
	"github.com/Merith-TK/utils/pkg/config"
	"github.com/Merith-TK/utils/pkg/debug"
	"github.com/Merith-TK/utils/pkg/driveutil"
)

func main() {
	flag.Parse()

	// Always enable debug for testing
	debug.SetDebug(true)

	fmt.Println("=== Testing all packages in pkg/ ===")

	testDebugPackage()
	testConfigPackage()
	testArchivePackage()
	testDriveutilPackage()

	fmt.Println("\n=== All package tests completed ===")

}

func testDebugPackage() {
	fmt.Println("\n--- Testing debug package ---")

	// Test title functionality
	debug.SetTitle("DEBUG_TEST")
	debug.Println("Testing with custom title")

	// Test title getter
	fmt.Printf("Current title: '%s'\n", debug.GetTitle())

	// Test debug state
	fmt.Printf("Debug enabled: %v\n", debug.GetDebug())

	// Test suicide mode (with very short timeout for testing)
	debug.Suicide(10) // 10 second timeout for testing

	// Reset title
	debug.ResetTitle()
	debug.Println("Testing with default title")

	// Test stacktrace if enabled
	if debug.GetStacktrace() {
		debug.Print("Stacktrace test")
	}
}

func testConfigPackage() {
	fmt.Println("\n--- Testing config package ---")

	// Test EnvKeyReplace
	replacements := map[string]string{
		"{USER}": "testuser",
		"{HOME}": "C:\\Users\\testuser",
		"{APP}":  "testapp",
	}

	testString := "User: {USER}, Home: {HOME}, App: {APP}"
	result := config.EnvKeyReplace(testString, replacements)
	fmt.Printf("EnvKeyReplace test:\n  Input: %s\n  Output: %s\n", testString, result)

	// Test EnvOverride
	testEnv := map[string]string{
		"TEST_VAR1": "value1",
		"TEST_VAR2": "value2",
	}

	config.EnvOverride(testEnv)
	fmt.Printf("EnvOverride test:\n  TEST_VAR1: %s\n  TEST_VAR2: %s\n",
		os.Getenv("TEST_VAR1"), os.Getenv("TEST_VAR2"))

	// Test TOML functionality with a temporary file
	testTOMLConfig()
}

func testTOMLConfig() {
	type TestConfig struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
		Debug   bool   `toml:"debug"`
	}

	// Create test config
	testCfg := TestConfig{
		Name:    "test-app",
		Version: "1.0.0",
		Debug:   true,
	}

	// Create temp file
	tempFile := filepath.Join(os.TempDir(), "test-config.toml")

	// Test SaveToml
	if err := config.SaveToml(tempFile, testCfg); err != nil {
		fmt.Printf("SaveToml error: %v\n", err)
		return
	}

	// Test LoadToml
	var loadedCfg TestConfig
	if err := config.LoadToml(&loadedCfg, tempFile); err != nil {
		fmt.Printf("LoadToml error: %v\n", err)
		return
	}

	fmt.Printf("TOML config test:\n  Original: %+v\n  Loaded: %+v\n", testCfg, loadedCfg)

	// Cleanup
	os.Remove(tempFile)
}

func testArchivePackage() {
	fmt.Println("\n--- Testing archive package ---")

	// Create a test directory structure
	testDir := filepath.Join(os.TempDir(), "archive-test")
	zipFile := filepath.Join(os.TempDir(), "test.zip")
	extractDir := filepath.Join(os.TempDir(), "extract-test")

	// Create test directory with some files
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)

	// Create test files
	testFiles := []string{
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(testDir, "file2.txt"),
		filepath.Join(testDir, "subdir", "file3.txt"),
	}

	for _, file := range testFiles {
		os.WriteFile(file, []byte("test content for "+filepath.Base(file)), 0644)
	}

	// Create zip file manually for testing
	if err := createTestZip(testDir, zipFile); err != nil {
		fmt.Printf("Failed to create test zip: %v\n", err)
		// Cleanup and return early
		os.RemoveAll(testDir)
		return
	}

	// Test Unzip
	if err := archive.Unzip(zipFile, extractDir); err != nil {
		fmt.Printf("Unzip error: %v\n", err)
	} else {
		fmt.Printf("Unzip successful: %s -> %s\n", zipFile, extractDir)

		// Verify extracted files
		extractedFiles := []string{
			filepath.Join(extractDir, "file1.txt"),
			filepath.Join(extractDir, "file2.txt"),
			filepath.Join(extractDir, "subdir", "file3.txt"),
		}

		for _, file := range extractedFiles {
			if _, err := os.Stat(file); err == nil {
				fmt.Printf("  ✓ Extracted: %s\n", strings.TrimPrefix(file, extractDir))
			} else {
				fmt.Printf("  ✗ Missing: %s\n", strings.TrimPrefix(file, extractDir))
			}
		}
	}

	// Cleanup
	os.RemoveAll(testDir)
	os.RemoveAll(extractDir)
	os.Remove(zipFile)
}

func createTestZip(sourceDir, zipFile string) error {
	// Create the zip file
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFileHandle.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

	// Walk through source directory and add files to zip
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path for zip entry
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// For directories, add trailing slash
		if info.IsDir() {
			relPath = relPath + "/"
		}

		// Create zip entry
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// If it's a directory, we're done (just created the entry)
		if info.IsDir() {
			return nil
		}

		// If it's a file, copy its contents
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})

	fmt.Printf("Created test zip: %s\n", zipFile)
	return err
}

func testDriveutilPackage() {
	fmt.Println("\n--- Testing driveutil package ---")

	// Create a drive store
	store := make(driveutil.DriveStore)

	// Test drive detection
	fmt.Println("Detecting drives...")
	store.DetectDrives(func(drive string, serial uint32) {
		fmt.Printf("  New drive detected: %s (Serial: %08X)\n", drive, serial)
	})

	// Test ListDrives
	fmt.Println("Listing all drives:")
	drives := driveutil.ListDrives()
	for _, drive := range drives {
		fmt.Printf("  Drive: %s, Label: %s, Serial: %08X, Type: %d\n",
			drive.Letter, drive.Label, drive.Serial, drive.Type)
	}

	// Test specific drive functions on found drives
	fmt.Println("Testing individual drive functions on found drives:")
	for _, drive := range drives {
		fmt.Printf("Testing drive: %s\n", drive.Letter)

		// Test GetVolumeSerialNumber
		if serial, err := driveutil.GetVolumeSerialNumber(drive.Letter); err == nil {
			fmt.Printf("  GetVolumeSerialNumber: %08X\n", serial)
			// Compare with ListDrives result
			if serial == drive.Serial {
				fmt.Printf("  ✓ Serial matches ListDrives result\n")
			} else {
				fmt.Printf("  ✗ Serial mismatch - ListDrives: %08X, GetVolumeSerialNumber: %08X\n", drive.Serial, serial)
			}
		} else {
			fmt.Printf("  GetVolumeSerialNumber: Error - %v\n", err)
		}

		// Test DriveExists
		if driveutil.DriveExists(drive.Letter) {
			fmt.Printf("  ✓ DriveExists confirms drive exists\n")
		} else {
			fmt.Printf("  ✗ DriveExists says drive doesn't exist (inconsistent)\n")
		}
	}

	// Test drive monitoring simulation (brief test)
	fmt.Println("Testing drive monitoring (3 second test)...")
	done := make(chan bool, 1)

	go func() {
		time.Sleep(3 * time.Second)
		done <- true
	}()

	go store.MonitorDrives(func(drive string, serial uint32) {
		fmt.Printf("  Monitor detected: %s (Serial: %08X)\n", drive, serial)
	}, 1*time.Second)

	<-done // Wait for monitoring test to complete
	fmt.Println("Drive monitoring test completed")
}
