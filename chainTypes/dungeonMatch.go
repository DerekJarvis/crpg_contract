package chainTypes

import "fmt"

type DungeonMatch struct {
	DungeonOwner        string
	DungeonName         string
	StartTime           int64
	EndTime             int64
	CharacterName       string
	Player              string
	DungeonMultiplier   float32
	CharacterMultiplier float32
	StartingPower       float32
	EndingPower         float32
	Reward              int
	PlayedTiles         []int
}

func (dungeonMatch DungeonMatch) Simulate(dungeon *Dungeon, tiles map[string]*Tile, character *Character, moves []int) (err error) {
	for _, v := range moves {
		currentTile := tiles[dungeon.Tiles[v].TileId]

		neededPower, reward := CalculateTile(character, currentTile, dungeonMatch.DungeonMultiplier, dungeonMatch.CharacterMultiplier)

		if character.Power-neededPower < 0 {
			// Deduct power, but tile was not defeated
			character.Power = 0
			break
		} else {
			dungeonMatch.Reward += reward

			// Update Power
			character.Power -= neededPower
			dungeon.Power -= neededPower

			// Add tile to list of played tiles
			dungeonMatch.PlayedTiles = append(dungeonMatch.PlayedTiles, v) // Tile needs ID and maybe dungeon tile? And maybe reference DT here?
		}
	}

	return nil
}

func CreateDungeonMatchIdParts(owner string, name string, startTime int64) []string {
	return []string{KEY_DUNGEON_MATCH, owner, name, fmt.Sprint(startTime)}
}

func (dungeonMatch DungeonMatch) GetIdParts() []string {
	return CreateDungeonMatchIdParts(string(dungeonMatch.DungeonOwner), dungeonMatch.DungeonName, dungeonMatch.StartTime)
}
