package chainTypes

import (
	"math/rand"
)

type Character struct {
	Id    string
	Name  string
	Owner WalletAddress
	ColorAttributes
	Power float32
	Gold  int
}

func CreateCharacterIdParts(id string) []string {
	return []string{KEY_CHARACTER, id}
}

func (character Character) GetIdParts() []string {
	return CreateCharacterIdParts(character.Id)
}

func (character Character) GetOwner() WalletAddress {
	return character.Owner
}

func NewCharacter(id string, name string, owner WalletAddress) *Character {
	c := new(Character)

	c.Id = id
	c.Name = name
	c.Owner = owner
	c.Red = rand.Intn(25) + 10
	c.Green = rand.Intn(25) + 10
	c.Blue = rand.Intn(25) + 10

	highestStat := GetHighestStat(&c.ColorAttributes)

	remainingPoints := CHARACTER_MAX_START_POINTS - c.Red - c.Green - c.Blue

	*highestStat += remainingPoints

	c.Power = CHARACTER_START_POWER

	return c
}

func (character Character) GetPower() (output int) {
	characterColors := ConvertColorsToArray(&character.ColorAttributes)

	for _, v := range characterColors {
		output += *v
	}

	return
}
