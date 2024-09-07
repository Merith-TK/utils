package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Merith-TK/utils/debug"
)

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "d", false, "Dry run")

	flag.Parse()
}

func main() {
	if dryRun {
		fmt.Println("Dry run enabled")
	}

	debug.Print("Args:", flag.Args())
	if flag.Args() != nil {
		debug.Print("Arguments provided")
		for _, arg := range flag.Args() {
			debug.Print("Arg:", arg)
			file, err := os.Open(arg)
			if err != nil {
				fmt.Println("Failed to open folder:", err)
				continue
			}

			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Println("Failed to get folder info:", err)
				continue
			}
			file.Close()

			sort(fileInfo)

		}

		return
	} else {
		debug.Print("No arguments provided")
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		return
	}

	for _, file := range files {
		sort(file)
	}
}
func sort(file fs.FileInfo) {
	if file.IsDir() {
		dir := file.Name()
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {

			cmd := exec.Command("git", "-C", dir, "remote", "get-url", "origin")
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("No origin found for", dir, "skipping...")
				return
			}

			url := strings.TrimSpace(string(output))
			base := strings.TrimSuffix(url, ".git")
			base = strings.TrimPrefix(base, "https://")
			base = strings.TrimPrefix(base, "http://")
			base = strings.TrimPrefix(base, "git@")
			base = strings.TrimSuffix(base, "/")
			base = filepath.ToSlash(base)
			debug.Print("URL:", url)
			debug.Print("Base:", base)

			parentDir := filepath.ToSlash(filepath.Dir(base))
			debug.Print("Parent:", parentDir)

			if !dryRun {
				fmt.Println("Moving", dir, "to", base)
				err = os.MkdirAll(parentDir, os.ModePerm)
				if err != nil {
					fmt.Println("Failed to create directory:", err)
					return
				}
				dest := filepath.Join(parentDir, dir)
				dest = filepath.ToSlash(dest)
				if _, err := os.Stat(dest); err == nil {
					fmt.Println("Destination already exists:", dest)
					return
				}
				err = os.Rename(dir, dest)
				if err != nil {
					fmt.Println("Failed to move directory:", err)
					return
				}
			} else {
				fmt.Println("Moving", dir, "to", base, "skipped (dry-run)")
			}

		}
	}
}
