package chainTypes

import (
	"math/rand"
)

type Tile struct {
	Id string
	TileType
	ObjectStatus
	ColorAttributes
	Owner WalletAddress
}

type TileType int

const (
	Ground TileType = iota
	Monster
	Loot
)

func CreateTileIdParts(id string) []string {
	return []string{KEY_TILE, id}
}

func (tile Tile) GetIdParts() []string {
	return CreateTileIdParts(tile.Id)
}

func (tile Tile) GetOwner() WalletAddress {
	return tile.Owner
}

// type TileDefinition struct {
// 	Collection string
// 	Id         int
// 	TileType
// }

func NewTile(id string, tileType TileType, owner WalletAddress) *Tile {

	t := new(Tile)
	t.TileType = tileType
	t.ObjectStatus = Available
	t.Magenta = rand.Intn(25) + 10
	t.Yellow = rand.Intn(25) + 10
	t.Cyan = rand.Intn(25) + 10
	t.Owner = owner
	t.Id = id

	return t
}

func (tile Tile) GetPower() (output int) {
	tileColors := ConvertColorsToArray(&tile.ColorAttributes)

	for _, v := range tileColors {
		output += *v
	}

	return
}
