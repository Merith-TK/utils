// Package main implements a Docker Compose wrapper that provides simplified commands
// for common Docker Compose operations while maintaining full compatibility.
//
// The dcomp utility offers convenient aliases for frequently used Docker Compose
// commands and passes through any unrecognized commands directly to docker compose.
//
// Available command aliases:
//   build  -> docker compose build
//   pull   -> docker compose pull
//   start  -> docker compose up -d --remove-orphans
//   stop   -> docker compose down
//   ps     -> docker compose ps
//   logs   -> docker compose logs
//
// Usage:
//   dcomp [command] [args...]
//
// Examples:
//   dcomp start          # Equivalent to: docker compose up -d --remove-orphans
//   dcomp stop           # Equivalent to: docker compose down
//   dcomp build web      # Equivalent to: docker compose build web
//   dcomp exec web bash  # Passes through: docker compose exec web bash
//
// Any command not in the alias list is passed directly to docker compose,
// ensuring full compatibility with all Docker Compose functionality.
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
