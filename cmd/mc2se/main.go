// Package main implements a Minecraft to Space Engineers schematic converter
// that transforms Litematica schematic files into Space Engineers blueprint format.
//
// The mc2se utility converts Minecraft structures saved in Litematica format
// into Space Engineers .sbc blueprint files, enabling cross-game structure transfer.
// It includes block mapping, color conversion, and grid size options.
//
// Features:
//   - Litematica schematic file loading and parsing
//   - Block-by-block conversion with material mapping
//   - Color conversion and block skin handling
//   - Small grid and large grid block support
//   - Custom definitions file support for block mapping
//   - Debug output integration for conversion tracking
//
// Usage:
//   mc2se -i [input-file] [-o output-file] [-d definitions-file] [-s]
//
// Flags:
//   -i    Input Litematica schematic file (required)
//   -o    Output Space Engineers blueprint file (default: output.sbc)
//   -d    Custom block definitions JSON file (default: definitions.json)
//   -s    Use small grid blocks instead of large grid
//
// Examples:
//   mc2se -i castle.litematic -o castle.sbc
//   mc2se -i house.litematic -o house.sbc -s
//   mc2se -i structure.litematic -d custom_blocks.json
//
// The converter processes each block in the schematic, maps it to appropriate
// Space Engineers blocks, and generates a complete .sbc blueprint file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Merith-TK/utils/pkg/debug"
	"github.com/elvis972602/go-litematica-tools/schematic"
)

var inputFile string
var outputFile string
var definitionsFile string
var smallGrid bool

func init() {
	flag.StringVar(&inputFile, "i", "", "input file name")
	flag.StringVar(&outputFile, "o", "output.sbc", "output file name")
	flag.StringVar(&definitionsFile, "d", "definitions.json", "custom definitions file")
	flag.BoolVar(&smallGrid, "s", false, "use small grid blocks")
	flag.Parse()
}

func main() {

	if inputFile == "" {
		log.Fatal("Input file not specified")
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

	debug.Print("Project metadata:", project.MetaData)
	debug.Print("Project version:", project.Version)
	debug.Print("Project minecraft data version:", project.MinecraftDataVersion)
	debug.Print("Project region name:", project.RegionName)
	if smallGrid {
		debug.Print("Using small grid blocks")
	} else {
		debug.Print("Using large grid blocks")
	}
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
				debug.SetTitle(blockState.Name)
				debug.Print("> ", blockType, color)
				debug.Print("> ", blockSkin)
				debug.Print("> ", x, y, z)
				debug.ResetTitle()

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
