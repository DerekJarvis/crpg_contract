package chainTypes

type PlayerVault struct {
	ExternalAddress WalletAddress
	Tiles           []string
	Packs           []string
}

func CreatePlayerVaultIdParts(externalAddress string) []string {
	return []string{KEY_PLAYER_VAULT, externalAddress}
}

func (playerVault PlayerVault) GetIdParts() []string {
	return CreatePlayerVaultIdParts(string(playerVault.ExternalAddress))
}

func NewPlayerVault(externalAddress string) (*PlayerVault, error) {
	playerVault := new(PlayerVault)

	playerVault.ExternalAddress = WalletAddress(externalAddress)
	playerVault.Tiles = make([]string, 0, 100)
	playerVault.Packs = make([]string, 0, 100)

	return playerVault, nil
}
