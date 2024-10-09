package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func startAutorun(drivePath string) {
	// Check if the drive path exists
	if _, err := os.Stat(drivePath); os.IsNotExist(err) {
		fmt.Printf("Drive path %s does not exist\n", drivePath)
		return
	}

	// Read the config file
	conf, err := setupConfig(drivePath + "/.autorun.toml")
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	if conf.Autorun == "" {
		fmt.Println("No autorun program specified")
		return
	}

	// Set the environment variables
	conf = setupEnvironment(conf)

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
		customEnv = isolatedEnv
		for _, value := range isolatedEnv {
			if _, err := os.Stat(value); os.IsNotExist(err) {
				if strings.HasPrefix(value, "C:") {
					fmt.Printf("Skipping directory creation for %s as it is on C: drive\n", value)
					continue
				}
				err := os.MkdirAll(value, os.ModePerm)
				if err != nil {
					fmt.Printf("Error creating directory %s: %s\n", value, err)
					return
				}
			}
		}
	} else {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				customEnv[parts[0]] = parts[1]
			}
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
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting autorun program: %s\n", err)
		return
	}

}
