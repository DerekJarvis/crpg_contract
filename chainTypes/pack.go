package chainTypes

import (
	"fmt"
	"math/rand"
)

type Pack struct {
	Id         string
	Collection string
	Owner      WalletAddress
}

// TODO: Define these a bit more, based on data. For now, keep it simple and random
type PackType struct {
	Name            string
	PackSize        int
	GroundReserved  int
	MonsterReserved int
	LootReserved    int
}

func CreatePackIdParts(id string) []string {
	return []string{KEY_PACK, id}
}

func (pack Pack) GetIdParts() []string {
	return CreatePackIdParts(pack.Id)
}

func (pack Pack) GetOwner() WalletAddress {
	return pack.Owner
}

func (packType PackType) Validate() (bool, string) {
	if packType.PackSize < packType.GroundReserved+
		packType.MonsterReserved+
		packType.LootReserved {

		return false, "Invalid pack definition: Reserved cards exceed pack size"
	} else {
		return true, ""
	}
}

func GetPackType(packType string) PackType {
	if packType != "Default" {
		panic("Only the 'Default' pack type is supported for now")
	}

	var packTypes = map[string]PackType{
		"Default": {Name: "Default", PackSize: PACK_SIZE, GroundReserved: 5, MonsterReserved: 3, LootReserved: 2},
	}

	returnPack := packTypes[packType]

	ok, message := returnPack.Validate()

	if ok {
		return returnPack
	} else {
		panic(message)
	}
}

func (pack Pack) OpenPack(baseId string) []*Tile {
	packType := GetPackType(pack.Collection)

	returnTiles := make([]*Tile, 0, packType.PackSize)

	currentPosition := 0

	for i := 0; i < packType.GroundReserved; i++ {
		returnTiles = append(returnTiles, (NewTile(fmt.Sprintf("%s_%d", baseId, currentPosition), Ground, pack.Owner)))
		currentPosition++
	}

	for i := 0; i < packType.MonsterReserved; i++ {
		returnTiles = append(returnTiles, (NewTile(fmt.Sprintf("%s_%d", baseId, currentPosition), Monster, pack.Owner)))
		currentPosition++
	}

	for i := 0; i < packType.LootReserved; i++ {
		returnTiles = append(returnTiles, (NewTile(fmt.Sprintf("%s_%d", baseId, currentPosition), Loot, pack.Owner)))
		currentPosition++
	}

	// For the remaining items, take a random number to assign the type, favoring the focus
	for i := currentPosition; i < packType.PackSize-1; i++ {
		typePicker := rand.Intn(3)

		var tileType TileType

		switch typePicker {
		case 1:
			tileType = Ground
		case 2:
			tileType = Monster
		case 3:
			tileType = Loot
		default:
			panic("Invalid randomization when generating a pack")
		}

		returnTiles = append(returnTiles, (NewTile(fmt.Sprintf("%s_%d", baseId, currentPosition), tileType, pack.Owner)))
	}

	return returnTiles
}
