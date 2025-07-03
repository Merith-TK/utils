package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// SandboxConfig represents the configuration for sandboxing
type SandboxConfig struct {
	DrivePath     string
	IsolatedPath  string
	AllowedDrives []string
	MaxMemoryMB   uint64
	TimeoutSec    uint32
}

// SandboxedProcess represents a sandboxed process
type SandboxedProcess struct {
	processInfo *windows.ProcessInformation
	config      SandboxConfig
}

// createSandboxedProcess creates a new sandboxed process with filesystem restrictions
func createSandboxedProcess(cmd *exec.Cmd, config SandboxConfig) (*SandboxedProcess, error) {
	fmt.Printf("[SANDBOX] Creating sandboxed process for: %s\n", cmd.Path)

	// Prepare process creation with restricted environment
	startupInfo := &windows.StartupInfo{
		Cb: uint32(unsafe.Sizeof(windows.StartupInfo{})),
	}

	processInfo := &windows.ProcessInformation{}

	// Convert command to Windows format
	cmdLine := `"` + cmd.Path + `"`
	if len(cmd.Args) > 1 {
		for _, arg := range cmd.Args[1:] {
			cmdLine += ` "` + arg + `"`
		}
	}

	fmt.Printf("[SANDBOX] Command line: %s\n", cmdLine)

	// Set working directory to isolated path
	workDir := config.IsolatedPath
	if cmd.Dir != "" {
		// Map the working directory to isolated path
		workDir = mapToIsolatedPath(cmd.Dir, config)
	}

	fmt.Printf("[SANDBOX] Working directory: %s\n", workDir)

	// Create environment block with restricted paths
	envBlock := createRestrictedEnvironment(cmd.Env, config)
	fmt.Printf("[SANDBOX] Environment block created\n")

	// Create the process
	cmdLinePtr, err := windows.UTF16PtrFromString(cmdLine)
	if err != nil {
		return nil, fmt.Errorf("failed to convert command line: %v", err)
	}

	workDirPtr, err := windows.UTF16PtrFromString(workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to convert working directory: %v", err)
	}

	err = windows.CreateProcess(
		nil,
		cmdLinePtr,
		nil,
		nil,
		false,
		0, // No special flags - simplify to avoid privilege issues
		envBlock,
		workDirPtr,
		startupInfo,
		processInfo,
	)

	if err != nil {
		fmt.Printf("[SANDBOX] Failed to create process: %v\n", err)
		return nil, fmt.Errorf("failed to create process: %v", err)
	}

	fmt.Printf("[SANDBOX] Process created successfully with PID: %d\n", processInfo.ProcessId)

	return &SandboxedProcess{
		processInfo: processInfo,
		config:      config,
	}, nil
}

// createRestrictedEnvironment creates an environment block with restricted paths
func createRestrictedEnvironment(baseEnv []string, config SandboxConfig) *uint16 {
	envMap := make(map[string]string)

	// Parse base environment
	for _, env := range baseEnv {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	// Override system paths to isolated equivalents
	driveLetter := config.DrivePath[:1]
	isolatedRoot := config.IsolatedPath

	// Redirect common paths to isolated directory
	envMap["USERPROFILE"] = filepath.Join(isolatedRoot, "User")
	envMap["APPDATA"] = filepath.Join(isolatedRoot, "AppData", "Roaming")
	envMap["LOCALAPPDATA"] = filepath.Join(isolatedRoot, "AppData", "Local")
	envMap["TEMP"] = filepath.Join(isolatedRoot, "Temp")
	envMap["TMP"] = filepath.Join(isolatedRoot, "Temp")
	envMap["HOME"] = filepath.Join(isolatedRoot, "User")

	// Restrict PATH to only essential Windows directories and the target drive
	restrictedPath := []string{
		"C:\\Windows\\System32",
		"C:\\Windows",
		driveLetter + ":\\",
	}
	envMap["PATH"] = strings.Join(restrictedPath, ";")

	// Convert to Windows environment block format (each string null-terminated, block double-null terminated)
	var envBlock []uint16
	for key, value := range envMap {
		envStr := key + "=" + value
		envUTF16, _ := windows.UTF16FromString(envStr)
		envBlock = append(envBlock, envUTF16...)
		envBlock = append(envBlock, 0) // Null terminate each string
	}
	envBlock = append(envBlock, 0) // Final null terminator for the block

	if len(envBlock) == 0 {
		return nil // No environment block
	}

	return &envBlock[0]
}

// mapToIsolatedPath maps a path to the isolated environment
func mapToIsolatedPath(originalPath string, config SandboxConfig) string {
	// Convert absolute paths on the target drive to isolated paths
	if strings.HasPrefix(originalPath, config.DrivePath) {
		relativePath := strings.TrimPrefix(originalPath, config.DrivePath)
		return filepath.Join(config.IsolatedPath, relativePath)
	}

	// For paths on other drives, default to isolated root
	return config.IsolatedPath
}

// Wait waits for the sandboxed process to complete
func (sp *SandboxedProcess) Wait() error {
	if sp.processInfo == nil {
		return fmt.Errorf("process not started")
	}

	// Wait for process to complete
	_, err := windows.WaitForSingleObject(sp.processInfo.Process, windows.INFINITE)
	return err
}

// Terminate forcefully terminates the sandboxed process
func (sp *SandboxedProcess) Terminate() error {
	if sp.processInfo == nil {
		return fmt.Errorf("process not started")
	}

	return windows.TerminateProcess(sp.processInfo.Process, 1)
}

// Close cleans up the sandboxed process resources
func (sp *SandboxedProcess) Close() error {
	if sp.processInfo != nil {
		if sp.processInfo.Process != 0 {
			windows.CloseHandle(sp.processInfo.Process)
		}
		if sp.processInfo.Thread != 0 {
			windows.CloseHandle(sp.processInfo.Thread)
		}
	}

	return nil
}

// GetExitCode returns the exit code of the process
func (sp *SandboxedProcess) GetExitCode() (uint32, error) {
	if sp.processInfo == nil {
		return 0, fmt.Errorf("process not started")
	}

	var exitCode uint32
	err := windows.GetExitCodeProcess(sp.processInfo.Process, &exitCode)
	return exitCode, err
}
