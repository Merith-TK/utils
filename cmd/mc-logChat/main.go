package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// check if there is an argument
	if len(os.Args) < 2 {
		fmt.Println("Error: No file specified")
		os.Exit(1)
	}
	arg := os.Args[1]

	file, err := ioutil.ReadFile(arg)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// convert file to string
	str := string(file)
	// split string into slice of strings by new line
	lines := strings.Split(str, "\n")

	chat := []string{}
	// loop through slice of strings
	for _, line := range lines {
		if strings.Contains(line, "[CHAT]") {
			chat = append(chat, line)
		}
	}

	// write chat to file
	err = ioutil.WriteFile("chat.txt", []byte(strings.Join(chat, "\n")), 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}
