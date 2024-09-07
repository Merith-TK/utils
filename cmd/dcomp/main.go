package main

import (
	"flag"
	"fmt"
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
	if commandBase, ok := commandList[command]; ok {
		execCmd(commandBase, args[1:]...)
	} else {
		log.Fatalf("Command '%s' not found", command)
	}

	fmt.Println("Command executed successfully")
}

func execCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run command: %s", err)
	}
}
