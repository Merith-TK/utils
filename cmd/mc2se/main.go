package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Merith-TK/utils/debug"
	"github.com/elvis972602/go-litematica-tools/schematic"
)

var inputFile string
var outputFile string
var definitionsFile string

func init() {
	flag.StringVar(&inputFile, "i", "", "input file name")
	flag.StringVar(&outputFile, "o", "output.sbc", "output file name")
	flag.StringVar(&definitionsFile, "d", "definitions.json", "custom definitions file")
	flag.Parse()
}

func main() {

	if inputFile == "" {
		if flag.Arg(0) != "" {
			inputFile = flag.Arg(0)
		} else {
			log.Fatal("Input file not specified")
		}
	}
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}
	defer file.Close()
	project, err := schematic.LoadFromFile(file)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Schematic loaded")
	}

	debug.Println("Project metadata:", project.MetaData)
	debug.Println("Project version:", project.Version)
	debug.Println("Project minecraft data version:", project.MinecraftDataVersion)
	debug.Println("Project region name:", project.RegionName)

	// Get the region size from the project.MetaData.EnclosingSize (vec3d)
	xSize := project.MetaData.EnclosingSize.X
	ySize := project.MetaData.EnclosingSize.Y
	zSize := project.MetaData.EnclosingSize.Z
	log.Println("Project region size:", xSize, ySize, zSize)

	loadDefinitions()
	convertColors()

	var blocklist string
	for x := 0; x < int(xSize); x++ {
		for y := 0; y < int(ySize); y++ {
			for z := 0; z < int(zSize); z++ {
				blockState := project.GetBlock(x, y, z)
				if blockState.Name == "minecraft:air" {
					continue
				}
				log.Println("Converting block:", blockState.Name)
				color, blockType, blockSkin := convertBlock(blockState)
				fmt.Println("\t> ", blockType, color)
				fmt.Println("\t> ", blockSkin)
				fmt.Println("\t> ", x, y, z)

				// todo: handle custom skins in conversion.go
				blocklist += writeBlock(blockType, color, []int{x, y, z}, project.RegionName, blockSkin)
			}
		}
	}

	// write the blueprint
	xmlOutput := xmlHeader + blocklist + xmlFooter

	// write the blueprint to file
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(xmlOutput)
	if err != nil {
		panic(err)
	}

}
