package chainTypes

type CharacterLookup struct {
	Id   string
	Name string
}

func CreateCharacterLookupIdParts(name string) []string {
	return []string{KEY_CHARACTER_LOOKUP, name}
}

func (characterLookup CharacterLookup) GetIdParts() []string {
	return CreateCharacterLookupIdParts(characterLookup.Name)
}
