package chainTypes

type DungeonListing struct {
	Name   string
	Owner  WalletAddress
	Status ObjectStatus
}

func CreateDungeonListingIdParts(owner string, name string) []string {
	return []string{KEY_DUNGEON_LISTING, owner, name}
}

func (dungeonListing DungeonListing) GetIdParts() []string {
	return CreateDungeonListingIdParts(string(dungeonListing.Owner), dungeonListing.Name)
}
