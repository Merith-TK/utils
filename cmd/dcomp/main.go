package main

import (
	"flag"
	"log"
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

	execCommand("docker", "compose", args...)

	// Rest of your code goes here
}
