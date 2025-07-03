package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Merith-TK/utils/pkg/config"
	"github.com/Merith-TK/utils/pkg/debug"
)

var conf Config

type Config struct {
	Autorun     string            `toml:"autorun,omitempty"`
	WorkDir     string            `toml:"workDir,omitempty"`
	Isolate     bool              `toml:"isolated,omitempty"`
	Environment map[string]string `toml:"environment,omitempty"`
}

func startAutorun(drivePath string) {
	log.Printf("[AUTORUN] Starting autorun check for drive: %s\n", drivePath)

	// Check if the drive path exists
	if _, err := os.Stat(drivePath); os.IsNotExist(err) {
		log.Printf("[AUTORUN] Drive path %s does not exist\n", drivePath)
		return
	}

	// Read the config file using pkg/config
	configPath := drivePath + "/.autorun.toml"
	log.Printf("[AUTORUN] Checking for config file: %s\n", configPath)
	err := config.LoadToml(&conf, configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[AUTORUN] Error reading config file: %s\n", err)
		} else {
			log.Printf("[AUTORUN] No config file found: %s\n", configPath)
		}
		return
	}

	if conf.Autorun == "" {
		log.Printf("[AUTORUN] No autorun program specified in config\n")
		return
	}

	log.Printf("[AUTORUN] Found autorun config: %s\n", conf.Autorun)

	// Set the environment variables using pkg/config
	conf = *setupEnvironment(&conf)

	if !filepath.IsAbs(conf.Autorun) {
		conf.Autorun = filepath.Join(drivePath, conf.Autorun)
	}
	if !filepath.IsAbs(conf.WorkDir) {
		conf.WorkDir = filepath.Join(drivePath, conf.WorkDir)
	}

	// start building the command
	cmd := exec.Command(conf.Autorun)
	cmd.Env = []string{}
	customEnv := map[string]string{}

	isolatedEnv := map[string]string{
		"HOME":              filepath.Join(drivePath, "/.isolated/User"),
		"USERPROFILE":       filepath.Join(drivePath, "/.isolated/User"),
		"APPDATA":           filepath.Join(drivePath, "/.isolated/AppData/Roaming"),
		"LOCALAPPDATA":      filepath.Join(drivePath, "/.isolated/AppData/Local"),
		"TEMP":              filepath.Join(drivePath, "/.isolated/Temp"),
		"TMP":               filepath.Join(drivePath, "/.isolated/Temp"),
		"SystemRoot":        "C:\\Windows",
		"ProgramFiles":      "C:\\Program Files",
		"ProgramFiles(x86)": "C:\\Program Files (x86)",
		"ProgramData":       "C:\\ProgramData",
	}

	if conf.Isolate {
		// Try advanced Windows sandboxing first, fall back to environment-only isolation
		log.Printf("[AUTORUN] Using advanced isolation mode for drive: %s", drivePath)

		// Create isolated directories
		isolatedRoot := filepath.Join(drivePath, ".isolated")
		for _, dir := range []string{
			filepath.Join(isolatedRoot, "User"),
			filepath.Join(isolatedRoot, "AppData", "Roaming"),
			filepath.Join(isolatedRoot, "AppData", "Local"),
			filepath.Join(isolatedRoot, "Temp"),
		} {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Printf("[AUTORUN] Error creating isolated directory %s: %s", dir, err)
				return
			}
		}

		// Prepare environment for execution
		customEnv = isolatedEnv
		for key, value := range conf.Environment {
			customEnv[key] = value
		}

		// Convert environment map to slice
		envSlice := []string{}
		for key, value := range customEnv {
			envSlice = append(envSlice, key+"="+value)
		}
		cmd.Env = envSlice
		cmd.Dir = conf.WorkDir
		if cmd.Dir == "" {
			cmd.Dir = isolatedRoot
		}

		// Try advanced sandboxing (requires admin privileges)
		sandboxConfig := SandboxConfig{
			DrivePath:     drivePath,
			IsolatedPath:  isolatedRoot,
			AllowedDrives: []string{drivePath[:1]}, // Only allow access to the target drive
			MaxMemoryMB:   512,                     // 512MB memory limit
			TimeoutSec:    300,                     // 5 minute timeout
		}

		// Attempt to create sandboxed process
		sandboxedProc, err := createSandboxedProcess(cmd, sandboxConfig)
		if err != nil {
			log.Printf("[AUTORUN] Advanced sandboxing failed (may require admin): %s", err)
			log.Printf("[AUTORUN] Falling back to environment-only isolation")

			// Fall back to regular process execution with isolated environment
			log.Printf("[AUTORUN] Starting command with environment isolation: %s (workdir: %s)", conf.Autorun, cmd.Dir)
			err = cmd.Start()
			if err != nil {
				log.Printf("[AUTORUN] Error starting autorun program: %s", err)
				return
			}
			log.Printf("[AUTORUN] Successfully started autorun program with environment isolation (PID: %d)", cmd.Process.Pid)
		} else {
			// Advanced sandboxing succeeded
			defer sandboxedProc.Close()

			log.Printf("[AUTORUN] Started sandboxed process for: %s", conf.Autorun)

			// Wait for completion
			err = sandboxedProc.Wait()
			if err != nil {
				log.Printf("[AUTORUN] Sandboxed process error: %s", err)
			}

			exitCode, _ := sandboxedProc.GetExitCode()
			log.Printf("[AUTORUN] Sandboxed process completed with exit code: %d", exitCode)
		}

	} else {
		// Use original environment-only isolation
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				customEnv[parts[0]] = parts[1]
			}
		}

		// Add custom environment variables
		for key, value := range conf.Environment {
			customEnv[key] = value
		}

		for key, value := range customEnv {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
		cmd.Dir = conf.WorkDir

		// Start the autorun program
		log.Printf("[AUTORUN] Starting command: %s (workdir: %s)", conf.Autorun, cmd.Dir)
		err = cmd.Start()
		if err != nil {
			log.Printf("[AUTORUN] Error starting autorun program: %s", err)
			return
		}

		log.Printf("[AUTORUN] Successfully started autorun program (PID: %d)", cmd.Process.Pid)
	}

}

// SetupEnvironment sets up environment variables and replaces placeholders in the config.
func setupEnvironment(conf *Config) *Config {
	debug.Print("config.SetupEnvironment called with config:", conf)
	// Variables for env replacement
	drivePath, _ := filepath.Abs("/")
	drivePath = filepath.ToSlash(drivePath)
	drivePath = strings.TrimSuffix(drivePath, "/")
	configEnvReplace := map[string]string{
		"{work}":  conf.WorkDir,
		"{drive}": drivePath,
	}
	debug.Print("Config environment replacements:", configEnvReplace)

	// Replace Normal Config options
	conf.Autorun = config.EnvKeyReplace(conf.Autorun, configEnvReplace)
	conf.WorkDir = config.EnvKeyReplace(conf.WorkDir, configEnvReplace)

	// Replace Environment Variables and set them
	for k, v := range conf.Environment {
		conf.Environment[k] = config.EnvKeyReplace(v, configEnvReplace)
	}
	config.EnvOverride(conf.Environment)

	return conf
}
