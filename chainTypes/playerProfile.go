package chainTypes

import "errors"

type WalletAddress string

type PlayerProfile struct {
	ExternalAddress WalletAddress
	Name            string
	Gold            int
	Characters      map[string]bool // Because go is lame and lacks basic types (like set)
	Dungeons        map[string]bool
}

func CreatePlayerProfileIdParts(externalAddress string) []string {
	return []string{KEY_PLAYER_PROFILE, externalAddress}
}

func (playerProfile PlayerProfile) GetIdParts() []string {
	return CreatePlayerProfileIdParts(string(playerProfile.ExternalAddress))
}

func GetSignatureMessage(address string, name string) string {
	return ""
}

func CheckSignatureMessage(signedMessage string) bool {
	return true
}

func NewPlayerProfile(name string, externalAddress string, signature string) (*PlayerProfile, error) {
	player := new(PlayerProfile)

	if !CheckSignatureMessage(signature) {
		return nil, errors.New("Signature is not valid")
	}

	player.Name = name
	player.ExternalAddress = WalletAddress(externalAddress)
	player.Gold = 1000
	player.Characters = make(map[string]bool, 10)
	player.Dungeons = make(map[string]bool, 10)

	return player, nil
}
