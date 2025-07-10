//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// CheckCGO checks if GCC is available for CGO builds
func CheckCGO() error {
	fmt.Println("Checking if GCC is available...")
	err := sh.Run("gcc", "--version")
	if err != nil {
		return fmt.Errorf("GCC not found in PATH. CGO requires gcc. Install with: scoop install gcc")
	}
	fmt.Println("GCC is available!")
	return nil
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("Running go mod tidy...")
	return sh.Run("go", "mod", "tidy")
}

// Build builds all commands in cmd/ folders
func Build() error {
	mg.Deps(CheckCGO, Tidy)
	
	fmt.Println("Building all commands...")
	
	// Create build directory if it doesn't exist
	if err := os.MkdirAll(".build", 0755); err != nil {
		return fmt.Errorf("failed to create .build directory: %w", err)
	}
	
	// Find all cmd folders
	cmdDirs, err := filepath.Glob("cmd/*")
	if err != nil {
		return fmt.Errorf("failed to find cmd directories: %w", err)
	}
	
	for _, cmdDir := range cmdDirs {
		if info, err := os.Stat(cmdDir); err != nil || !info.IsDir() {
			continue
		}
		
		fmt.Printf("Building %s\n", cmdDir)
		
		// Check for .buildargs file
		var args []string
		buildargsFile := filepath.Join(cmdDir, ".buildargs")
		if _, err := os.Stat(buildargsFile); err == nil {
			content, err := os.ReadFile(buildargsFile)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", buildargsFile, err)
			}
			
			argsStr := strings.TrimSpace(string(content))
			if argsStr != "" {
				fmt.Printf("Args: %s\n", argsStr)
				args = strings.Fields(argsStr)
			}
		}
		
		// Build the command
		buildArgs := []string{"build", "-o", ".build/"}
		buildArgs = append(buildArgs, args...)
		buildArgs = append(buildArgs, "./"+cmdDir)
		
		if err := sh.Run("go", buildArgs...); err != nil {
			return fmt.Errorf("failed to build %s: %w", cmdDir, err)
		}
	}
	
	return nil
}

// Install installs all commands in cmd/ folders
func Install() error {
	mg.Deps(CheckCGO, Tidy)
	
	fmt.Println("Installing all commands...")
	
	// Find all cmd folders
	cmdDirs, err := filepath.Glob("cmd/*")
	if err != nil {
		return fmt.Errorf("failed to find cmd directories: %w", err)
	}
	
	for _, cmdDir := range cmdDirs {
		if info, err := os.Stat(cmdDir); err != nil || !info.IsDir() {
			continue
		}
		
		fmt.Printf("Installing %s\n", cmdDir)
		
		// Check for .buildargs file
		var args []string
		buildargsFile := filepath.Join(cmdDir, ".buildargs")
		if _, err := os.Stat(buildargsFile); err == nil {
			content, err := os.ReadFile(buildargsFile)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", buildargsFile, err)
			}
			
			argsStr := strings.TrimSpace(string(content))
			if argsStr != "" {
				fmt.Printf("Args: %s\n", argsStr)
				args = strings.Fields(argsStr)
			}
		}
		
		// Install the command
		installArgs := []string{"install"}
		installArgs = append(installArgs, args...)
		installArgs = append(installArgs, "./"+cmdDir)
		
		if err := sh.Run("go", installArgs...); err != nil {
			return fmt.Errorf("failed to install %s: %w", cmdDir, err)
		}
	}
	
	return nil
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	return sh.Rm(".build")
}

// Help displays available targets
func Help() {
	fmt.Println("Available targets:")
	fmt.Println("  build    - Build all commands")
	fmt.Println("  install  - Install all commands")
	fmt.Println("  tidy     - Run go mod tidy")
	fmt.Println("  clean    - Clean build artifacts")
	fmt.Println("  checkCGO - Check if GCC is available")
	fmt.Println("  help     - Show this help")
}

// Default target
var Default = Build
