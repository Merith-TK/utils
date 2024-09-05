package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	// convert the environment map to a slice
	env := make([]string, 0, len(conf.Environment))

	if !filepath.IsAbs(conf.Autorun) {
		conf.Autorun = filepath.Join(drivePath, conf.Autorun)
	}
	if !filepath.IsAbs(conf.WorkDir) {
		conf.WorkDir = filepath.Join(drivePath, conf.WorkDir)
	}

	// start building the command
	cmd := exec.Command(conf.Autorun)

	if conf.Isolate {
		cmd.Env = env
		cmd.Env = append(cmd.Env, "TEMP="+os.Getenv("TEMP"))
	} else {
		cmd.Env = append(os.Environ(), env...)
	}

	cmd.Dir = conf.WorkDir

	// Start the autorun program
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting autorun program: %s\n", err)
		return
	}

}
