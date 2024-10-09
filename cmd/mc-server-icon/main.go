package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"

	_ "image/gif"
	"image/png"

	_ "image/jpeg"

	_ "image/png"
	"os"

	"golang.org/x/image/draw"
)

func init() {
	flag.String("icon", "", "Path to the icon file")
	flag.String("output", "server-icon.png", "Path to the output file")
	flag.Bool("help", false, "Show help")
}

func main() {
	flag.Parse()

	icon := flag.Lookup("icon").Value.String()
	output := flag.Lookup("output").Value.String()
	// move := flag.Lookup("move").Value.String() == "true"
	help := flag.Lookup("help").Value.String() == "true"
	if help || icon == "" {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(icon)
	if err != nil {
		fmt.Printf("Error opening icon file: %v\n", err)
		return
	}
	defer file.Close()

	_, _, err = image.Decode(file)
	if err != nil {
		fmt.Printf("Error: %s is not a valid image file\n", icon)
		return
	}

	file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return
	}

	resizedImg := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	outFile, err := os.Create(output)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outFile.Close()

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, resizedImg); err != nil {
		fmt.Printf("Error encoding image to PNG: %v\n", err)
		return
	}
	if _, err := outFile.Write(buffer.Bytes()); err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

}
