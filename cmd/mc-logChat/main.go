package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
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

	str := string(file)
	lines := strings.Split(str, "\n")

	chat := []string{}
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
