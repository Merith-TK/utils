package main

import (
	"strconv"
	"strings"
)

var xmlHeader = `<?xml version="1.0"?>
<Definitions xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<ShipBlueprints>
		<ShipBlueprint xsi:type="MyObjectBuilder_ShipBlueprintDefinition">
			<Id Type="MyObjectBuilder_ShipBlueprintDefinition" Subtype="{NAME}" />
			<DisplayName>{NAME}</DisplayName>
			<CubeGrids>
				<CubeGrid>
					<SubtypeName />
					<EntityId>0</EntityId>
					<PersistentFlags>CastShadows InScene</PersistentFlags>
					<PositionAndOrientation>
						<Position x="0" y="0" z="0" />
						<Forward x="0" y="0" z="0" />
						<Up x="0" y="0" z="0" />
						<Orientation>
							<X>0</X>
							<Y>0</Y>
							<Z>0</Z>
							<W>0</W>
						</Orientation>
					</PositionAndOrientation>
					<LocalPositionAndOrientation xsi:nil="true" />
					<GridSizeEnum>{GRID}</GridSizeEnum>
					<CubeBlocks>
`
var xmlFooter = `
					</CubeBlocks>
					<LinearVelocity x="0" y="0" z="0" />
					<AngularVelocity x="0" y="0" z="0" />
					<DisplayName>{NAME}</DisplayName>
					<DestructibleBlocks>true</DestructibleBlocks>
					<IsRespawnGrid>false</IsRespawnGrid>
					<LocalCoordSys>0</LocalCoordSys>
					<TargetingTargets />
				</CubeGrid>
			</CubeGrids>
			<EnvironmentType>None</EnvironmentType>
			<WorkshopId>0</WorkshopId>
			<OwnerSteamId>0</OwnerSteamId>
			<Points>0</Points>
		</ShipBlueprint>
	</ShipBlueprints>
</Definitions>
`

// <SubtypeName>LargeBlockArmorBlock</SubtypeName>
// <SubtypeName>LargeHeavyBlockArmorBlock</SubtypeName>
var blockTemplate = `
						<MyObjectBuilder_CubeBlock xsi:type="MyObjectBuilder_CubeBlock">
							<SubtypeName>{BLOCKTYPE}</SubtypeName>
							<ColorMaskHSV x="{COLOR1}" y="{COLOR2}" z="{COLOR3}" />
							<SkinSubtypeId>{SKIN}</SkinSubtypeId>
							<Min x="{POSX}" y="{POSY}" z="{POSZ}" />
							<BuiltBy>0</BuiltBy>
						</MyObjectBuilder_CubeBlock>
`

// blocktype, HSV color as three floats
func writeBlock(blockType string, color []float64, pos []int, bpName string, skin string) string {
	newBlock := blockTemplate

	// TODO: There has to be a better way to do this?
	if smallGrid && strings.HasPrefix(blockType, "Large") {
		blockType = strings.Replace(blockType, "Large", "Small", 1)
	}
	newBlock = strings.Replace(newBlock, "{BLOCKTYPE}", blockType, 1)
	newBlock = strings.Replace(newBlock, "{COLOR1}", strconv.FormatFloat(color[0], 'f', 6, 64), 1)
	newBlock = strings.Replace(newBlock, "{COLOR2}", strconv.FormatFloat(color[1], 'f', 6, 64), 1)
	newBlock = strings.Replace(newBlock, "{COLOR3}", strconv.FormatFloat(color[2], 'f', 6, 64), 1)
	newBlock = strings.Replace(newBlock, "{POSX}", strconv.Itoa(pos[0]), 1)
	newBlock = strings.Replace(newBlock, "{POSY}", strconv.Itoa(pos[1]), 1)
	newBlock = strings.Replace(newBlock, "{POSZ}", strconv.Itoa(pos[2]), 1)

	newBlock = strings.Replace(newBlock, "{SKIN}", skin, 1)

	xmlFooter = strings.Replace(xmlFooter, "{NAME}", bpName, -1)
	xmlHeader = strings.Replace(xmlHeader, "{NAME}", bpName, -1)

	if smallGrid {
		xmlHeader = strings.Replace(xmlHeader, "{GRID}", "Small", -1)
	} else {
		xmlHeader = strings.Replace(xmlHeader, "{GRID}", "Large", -1)
	}
	return newBlock
}
