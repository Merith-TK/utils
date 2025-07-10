package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("No command provided")
	}

	command := args[0]
	if _, ok := commandList[command]; !ok {
		commandHandler(command, args...)
		return
	}

	// Build the full command with "docker compose" + args
	fullArgs := append([]string{"compose"}, args...)
	execCommand("docker", fullArgs...)
}

func commandHandler(command string, args ...string) {
	log.Printf("Unknown command: %s", command)
	log.Printf("Available commands: %v", getAvailableCommands())
	os.Exit(1)
}

func execCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		log.Printf("Command failed: %v", err)
		os.Exit(1)
	}
}

func getAvailableCommands() []string {
	var commands []string
	for cmd := range commandList {
		commands = append(commands, cmd)
	}
	return commands
}
