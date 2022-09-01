package chainTypes

type DungeonTile struct {
	X      int
	Y      int
	TileId string
}

type Dungeon struct {
	Name       string
	Owner      WalletAddress
	Width      int
	Height     int
	Power      float32
	Difficulty int
	Gold       int
	Tiles      []DungeonTile
}

func CreateDungeonIdParts(owner string, name string) []string {
	return []string{KEY_DUNGEON, owner, name}
}

func (dungeon Dungeon) GetIdParts() []string {
	return CreateDungeonIdParts(string(dungeon.Owner), dungeon.Name)
}

func (dungeon Dungeon) GetOwner() WalletAddress {
	return dungeon.Owner
}

func GetCoordinates(position int, width int, height int) (x int, y int) {
	x = position % width
	y = position / width

	if y > height {
		y = height
	}

	return
}

func GetPosition(x, y, width, height int) (position int) {
	position = y*width + x

	return
}

func CalculateTile(character *Character, tile *Tile, characterMultiplier float32, dungeonMultiplier float32) (neededPower float32, reward int) {
	colorDifference := float32(SmashColors(&character.ColorAttributes, &tile.ColorAttributes)) *
		dungeonMultiplier * characterMultiplier

	// Calculate Reward
	switch tile.TileType {
	case Ground:
		neededPower = colorDifference / 255
		reward = 10
	case Monster:
		neededPower = colorDifference / 25.5
		reward = 10 * int(neededPower)
	case Loot:
		neededPower = colorDifference / 25.5
		reward = 10 * int(neededPower)
	}

	return
}
