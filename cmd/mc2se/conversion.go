package main

import (
	"encoding/json"
	"os"

	"github.com/Merith-TK/utils/debug"
	"github.com/elvis972602/go-litematica-tools/schematic"
	"github.com/lucasb-eyer/go-colorful"
)

var validColorsHex = map[string]string{
	// known minecraft colors
	"pink":       "#f38baf",
	"magenta":    "#c74ebd",
	"purple":     "#8932b8",
	"blue":       "#3c44aa",
	"light_blue": "#3ab3da",
	"cyan":       "#169c9c",
	"green":      "#5e7c16",
	"lime":       "#80c71f",
	"yellow":     "#f7e9a3",
	"orange":     "#f9801d",
	"red":        "#b02e26",
	"brown":      "#835432",
	"black:":     "#1d1d21",
	"gray":       "#474f52",
	"light_gray": "#9d9d97",
	"white":      "#f9fffe",
}

var validBlockTypes = map[string]string{
	"minecraft:wool":     "LargeBlockArmorBlock",
	"minecraft:concrete": "LargeHeavyBlockArmorBlock",
}

// var string: value [int, int, int]
var validColorsHsv = map[string][]float64{
	// convert color name to HSV colors as three floats

}

func convertColors() {
	// convert hex colors to HSV colors that
	// are compatible with Space Engineers
	for _, colorHex := range validColorsHex {
		c, err := colorful.Hex(colorHex)
		if err != nil {
			continue
		}
		h, s, v := c.Hsv()
		h = h / 360
		s = s - 0.8
		v = v - 0.45
		validColorsHsv[colorHex] = []float64{h, s, v}
	}
}

// convertBlock converts a block state into a color and a block type
// returns the color and the block type
func convertBlock(block schematic.BlockState) ([]float64, string, string) {
	// first check if the block is a known block
	for _, blockDef := range blockDefinitions {
		if blockDef.Block == block.Name {
			return validColorsHsv[blockDef.Color], blockDef.BlockType, blockDef.BlockSkin
		}
	}
	// if the block is not known, return a default block
	return validColorsHsv["#f9fffe"], "LargeBlockArmorBlock", "Concrete_Armor"
}

var blockDefinitions []blockDefinition

type blockDefinition struct {
	Block     string `json:"block"`
	Color     string `json:"color"`
	BlockType string `json:"blockType,omitempty"`
	BlockSkin string `json:"blockSkin,omitempty"`
}

func loadDefaultDefinitions() {
	for i := 0; i < 2; i++ {
		var mcType string
		if i == 0 {
			mcType = "minecraft:wool"
		} else {
			mcType = "minecraft:concrete"
		}
		for color, colorHex := range validColorsHex {
			var blockDef []blockDefinition
			blockDef = append(blockDef, blockDefinition{Block: mcType + "_" + color, Color: colorHex, BlockType: validBlockTypes[mcType], BlockSkin: "Concrete_Armor"})
			blockDefinitions = append(blockDefinitions, blockDef...)
			debug.Println("Generated Default Block Definition:\n", blockDef)
		}
	}
}

func loadDefinitions() {
	debug.Println("Loading definitions")
	if _, err := os.Stat(definitionsFile); os.IsNotExist(err) {
		// generate default definitions for stone and oak_planks
		exampleDefinitions := []blockDefinition{
			{Block: "minecraft:stone", Color: "#bfbfbf", BlockType: "LargeHeavyBlockArmorBlock"},
			{Block: "minecraft:oak_planks", Color: "#835432", BlockType: "LargeHeavyBlockArmorBlock"},
		}
		jsonData, err := json.MarshalIndent(exampleDefinitions, "", "    ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(definitionsFile, jsonData, 0644)
		if err != nil {
			panic(err)
		}
	}
	definitionsFile, err := os.ReadFile(definitionsFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(definitionsFile, &blockDefinitions)
	if err != nil {
		panic(err)
	}
	debug.Println("Loaded Custom Block Definitions:\n", blockDefinitions)

	loadDefaultDefinitions()
}
