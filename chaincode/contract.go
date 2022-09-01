package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"jarvispowered.com/color_rpg/chain/chainTypes"
)

type SmartContract struct {
	contractapi.Contract
}

func GetChainKeyByParts(ctx contractapi.TransactionContextInterface,
	idParts []string) (string, error) {

	return ctx.GetStub().CreateCompositeKey(idParts[0], idParts[1:])
}

func GetChainKeyByObject(ctx contractapi.TransactionContextInterface,
	chainObject chainTypes.ChainObject) (string, error) {

	chainKeyParts := chainObject.GetIdParts()
	return ctx.GetStub().CreateCompositeKey(chainKeyParts[0], chainKeyParts[1:])
}

func DeleteChainObject(ctx contractapi.TransactionContextInterface,
	chainObject chainTypes.ChainObject) error {

	chainKey, err := GetChainKeyByObject(ctx, chainObject)

	if err != nil {
		return nil
	}

	err = ctx.GetStub().DelState(chainKey)

	return err
}

func WriteChainObject(ctx contractapi.TransactionContextInterface,
	chainObject chainTypes.ChainObject) error {

	chainKey, err := GetChainKeyByObject(ctx, chainObject)

	if err != nil {
		return nil
	}

	chainPayload, err := json.Marshal(chainObject)

	if err != nil {
		return nil
	}

	err = ctx.GetStub().PutState(chainKey, chainPayload)

	return err
}

func ReadChainObjectByParts(ctx contractapi.TransactionContextInterface,
	idParts []string, chainObject chainTypes.ChainObject) (string, error) {

	chainKey, err := GetChainKeyByParts(ctx, idParts)

	if err != nil {
		return "", err
	}

	chainState, err := ctx.GetStub().GetState(chainKey)

	if err != nil {
		return "", err
	}

	// Marshal/Unmarshal will deserialize to the base type, and not just the interface being referred to
	err = json.Unmarshal(chainState, chainObject)

	if err != nil {
		return "", fmt.Errorf("Error: %s\nChainObject:\n%s\nChainKey:\n%s\nChainState:\n%s", err.Error(), chainObject, chainKey, string(chainState))
	} else {
		return chainKey, err
	}
}

func ReadChainObjectByRef(ctx contractapi.TransactionContextInterface,
	chainObject chainTypes.ChainObject) (string, error) {

	idParts := chainObject.GetIdParts()

	returnError := fmt.Errorf("Cannot read chain object by reference:\nObject:%s\nId Parts:\n%s", idParts, chainObject)

	hasError := false

	idPartCount := len(idParts)

	if idPartCount == 0 {
		// If there are no parts
		hasError = true
	} else if idPartCount > 1 && idParts[0] == "" {
		// If more than one ID part and first is emtpy
		hasError = true
	}

	if hasError {
		return "", returnError
	} else {
		return ReadChainObjectByParts(ctx, idParts, chainObject)
	}
}

// TODO: Add signature checks on all necessary functions

func getPlayerProfile(ctx contractapi.TransactionContextInterface,
	externalAddress string) (*chainTypes.PlayerProfile, string, error) {

	var playerProfile *chainTypes.PlayerProfile = new(chainTypes.PlayerProfile)

	chainKey, err := ReadChainObjectByParts(ctx, chainTypes.CreatePlayerProfileIdParts(externalAddress), playerProfile)

	if err != nil {
		return nil, "", err
	} else {
		return playerProfile, chainKey, nil
	}

}

func (s *SmartContract) GetPlayerProfile(ctx contractapi.TransactionContextInterface,
	externalAddress string) (*chainTypes.PlayerProfile, error) {

	playerProfile, _, err := getPlayerProfile(ctx, externalAddress)

	if err != nil {
		return nil, err
	} else {
		return playerProfile, nil
	}

}

func getPlayerVault(ctx contractapi.TransactionContextInterface,
	externalAddress string) (*chainTypes.PlayerVault, string, error) {

	var playerVault *chainTypes.PlayerVault = new(chainTypes.PlayerVault)

	chainKey, err := ReadChainObjectByParts(ctx, chainTypes.CreatePlayerVaultIdParts(externalAddress), playerVault)

	if err != nil {
		return nil, "", err
	} else {
		return playerVault, chainKey, nil
	}

}

func (s *SmartContract) GetPlayerVault(ctx contractapi.TransactionContextInterface,
	externalAddress string) (*chainTypes.PlayerVault, error) {

	playerVault, _, err := getPlayerVault(ctx, externalAddress)

	if err != nil {
		return nil, err
	} else {
		return playerVault, nil
	}

}

func getPacks(ctx contractapi.TransactionContextInterface,
	packIds []string) ([]*chainTypes.Pack, []string, error) {
	returnPacks := make([]*chainTypes.Pack, 0, len(packIds))
	returnKeys := make([]string, 0, len(packIds))

	for _, packId := range packIds {
		var currentPack *chainTypes.Pack = new(chainTypes.Pack)

		chainKey, err := ReadChainObjectByParts(ctx, chainTypes.CreatePackIdParts(packId), currentPack)

		if err != nil {
			return nil, nil, err
		} else {
			returnPacks = append(returnPacks, currentPack)
			returnKeys = append(returnKeys, chainKey)
		}
	}

	return returnPacks, returnKeys, nil
}

func (s *SmartContract) GetPacks(ctx contractapi.TransactionContextInterface,
	packIds []string) ([]*chainTypes.Pack, error) {

	returnPacks, _, err := getPacks(ctx, packIds)

	if err != nil {
		return nil, err
	} else {
		return returnPacks, err
	}
}

func getTiles(ctx contractapi.TransactionContextInterface,
	tileIds []string) ([]*chainTypes.Tile, []string, error) {
	returnTiles := make([]*chainTypes.Tile, 0, len(tileIds))
	returnKeys := make([]string, 0, len(tileIds))

	for _, tileId := range tileIds {
		var currentTile *chainTypes.Tile = new(chainTypes.Tile)

		chainKey, err := ReadChainObjectByParts(ctx, chainTypes.CreateTileIdParts(tileId), currentTile)

		if err != nil {
			return nil, nil, err
		} else {
			returnTiles = append(returnTiles, currentTile)
			returnKeys = append(returnKeys, chainKey)
		}
	}

	return returnTiles, returnKeys, nil
}

func (s *SmartContract) GetTiles(ctx contractapi.TransactionContextInterface,
	tileIds []string) ([]*chainTypes.Tile, error) {

	returnPacks, _, err := getTiles(ctx, tileIds)

	if err != nil {
		return nil, err
	} else {
		return returnPacks, err
	}
}

func (s *SmartContract) GetPlayerTiles(ctx contractapi.TransactionContextInterface,
	ethWallet string) ([]*chainTypes.Tile, error) {

	// Get player vault, if it fails, throw the error
	playerVault, _, err := getPlayerVault(ctx, ethWallet)

	if err != nil {
		return nil, err
	}

	// Get tiles for player
	returnTiles, _, err := getTiles(ctx, playerVault.Tiles)

	if err != nil {
		return nil, err
	} else {
		return returnTiles, err
	}
}

func (s *SmartContract) NewPlayerProfile(ctx contractapi.TransactionContextInterface,
	name string, externalAddress string, signature string) (*chainTypes.PlayerProfile, error) {

	// TODO: Check that the externalAddress is valid

	_, _, err := getPlayerProfile(ctx, externalAddress)

	if err == nil {
		return nil, errors.New(fmt.Sprintf("Player Profile already exists for %s\n%s", externalAddress, err.Error()))
	}

	_, _, err = getPlayerVault(ctx, externalAddress)

	if err == nil {
		return nil, errors.New(fmt.Sprintf("Player Vault already exists for %s\n%s", externalAddress, err.Error()))
	}

	// Construct the player we want to write
	newPlayerProfile, err := chainTypes.NewPlayerProfile(name, externalAddress, signature)

	if err != nil {
		return nil, err
	}

	newPlayerVault, err := chainTypes.NewPlayerVault(externalAddress)

	if err != nil {
		return nil, err
	}

	writeObjects := []chainTypes.ChainObject{newPlayerProfile, newPlayerVault}

	// We've confirmed both the profile and wallet do not exist, now we write to chain
	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return nil, err
		}
	}

	return newPlayerProfile, nil
}

func (s *SmartContract) GetCharacterByName(ctx contractapi.TransactionContextInterface,
	name string) (*chainTypes.Character, error) {

	characterId, err := s.GetCharacterIdByName(ctx, name)

	if err != nil {
		return nil, err
	}

	character := &chainTypes.Character{Id: characterId}

	_, err = ReadChainObjectByRef(ctx, character)

	if err != nil {
		return nil, err
	}

	return character, nil
}

func (s *SmartContract) GetCharacterIdByName(ctx contractapi.TransactionContextInterface,
	name string) (string, error) {

	characterLookup := &chainTypes.CharacterLookup{Name: name}

	_, err := ReadChainObjectByRef(ctx, characterLookup)

	if err != nil {
		return "", err
	}

	return characterLookup.Id, nil
}

func (s *SmartContract) NewCharacter(ctx contractapi.TransactionContextInterface,
	name string, externalAddress string, signature string) (*chainTypes.Character, error) {

	_ = signature // Ignore for now until we have signatures in place

	profileChainState, _, err := getPlayerProfile(ctx, externalAddress)

	// Error fetching the player's profile, return it
	if err != nil {
		return nil, err
	}

	// Make sure the character name doesn't exist already
	_, err = s.GetCharacterIdByName(ctx, name)

	if err == nil {
		return nil, errors.New(fmt.Sprintf("Character with name %s exists already", name))
	}

	txId := ctx.GetStub().GetTxID()

	// Key the characters based on the txid to avoid key collissions (either on player-based IDs or a global incrementor)
	newCharacter := chainTypes.NewCharacter(txId, name, chainTypes.WalletAddress(externalAddress))

	newCharacterLookup := chainTypes.CharacterLookup{Id: txId, Name: name}

	// Now we need to update the profile to include the character and confirm it serializes properly
	profileChainState.Characters[newCharacter.Name] = true

	// We've confirmed all the things, now let's write
	writeObjects := []chainTypes.ChainObject{profileChainState, newCharacter, newCharacterLookup}

	// We've confirmed both the profile and wallet do not exist, now we write to chain
	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return nil, err
		}
	}

	return newCharacter, nil
}

// Will deduct the appropriate amount of gold from the player, create packs, and assign to player
func (s *SmartContract) BuyPacks(ctx contractapi.TransactionContextInterface,
	owner string, quantity int, signature string) (*chainTypes.PlayerVault, error) {

	if quantity <= 0 {
		return nil, errors.New(fmt.Sprintf("Invalid quantity of packs: %d", quantity))
	}

	_ = signature // Ignore for now until we have signatures in place

	profileChainState, _, err := getPlayerProfile(ctx, owner)

	// Error fetching the player's profile, return error
	if err != nil {
		return nil, err
	}

	vaultChainState, _, err := getPlayerVault(ctx, owner)

	// Error fetching the player's profile, return error
	if err != nil {
		return nil, err
	}

	// Then confirm they have enough gold
	totalPackCost := quantity * chainTypes.PACK_COST
	if profileChainState.Gold < totalPackCost {
		return nil, errors.New(fmt.Sprintf("Player has insufficient gold (%d) to purchase %d packs for %d", profileChainState.Gold, quantity, totalPackCost))
	}

	// Then deduct the gold
	profileChainState.Gold -= totalPackCost

	txId := ctx.GetStub().GetTxID()

	writeObjects := make([]chainTypes.ChainObject, 0, quantity+2) // # of packs + vault + profile

	for i := 0; i < quantity; i++ {
		packId := fmt.Sprintf("%s_%d", txId, i)
		newPack := chainTypes.Pack{Id: packId, Collection: "Default", Owner: chainTypes.WalletAddress(owner)}

		writeObjects = append(writeObjects, newPack)

		vaultChainState.Packs = append(vaultChainState.Packs, newPack.Id)
	}

	writeObjects = append(writeObjects, vaultChainState)

	writeObjects = append(writeObjects, profileChainState)

	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return nil, err
		}
	}

	return vaultChainState, nil
}

func (s *SmartContract) OpenPacks(ctx contractapi.TransactionContextInterface,
	owner string, packIds []string, signature string) (*chainTypes.PlayerVault, error) {

	_ = signature // Ignore for now until we have signatures in place

	// Let's get the user first, in case it's invalid that will be quicker

	ownerWallet := chainTypes.WalletAddress(owner)

	vault, _, err := getPlayerVault(ctx, owner)

	// Initialize random with the timestamp
	// TODO: This is a risk for decentralization, because timestamp can be manipulated by the client
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	rand.Seed(timestamp.Seconds)

	txId := ctx.GetStub().GetTxID()

	if err != nil {
		return nil, err
	}

	packs, _, err := getPacks(ctx, packIds)

	if err != nil {
		return nil, err
	}

	writeObjects := make([]chainTypes.ChainObject, 0, len(packs)*chainTypes.PACK_SIZE+1) // # of tiles in packs + profile
	deleteObjects := make([]chainTypes.ChainObject, 0, len(packs))                       // # of packs

	for i, pack := range packs {
		if pack.Owner != ownerWallet {
			return nil, errors.New(fmt.Sprintf("Attempted to open a pack not owned by the user: %s", pack.Id))
		}

		// Remove from the vault to save another loop
		success := false
		vault.Packs, success = removeStringItem(vault.Packs, pack.Id)

		if !success {
			panic(fmt.Sprintf("Failed to remove a pack (%s) from a vault (%s) that should have contained it!", pack.Id, owner))
		}

		// Open the pack and add to the vault
		newTiles := pack.OpenPack(txId + "_" + strconv.Itoa(i))

		for _, tile := range newTiles {
			writeObjects = append(writeObjects, tile)

			// Append to vault
			vault.Tiles = append(vault.Tiles, tile.Id)
		}

		deleteObjects = append(deleteObjects, pack)

	}

	writeObjects = append(writeObjects, vault)

	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return nil, err
		}
	}

	for _, o := range deleteObjects {
		err = DeleteChainObject(ctx, o)
		if err != nil {
			return nil, err
		}
	}

	return vault, nil

}

func (s *SmartContract) MakeDungeon(ctx contractapi.TransactionContextInterface,
	owner string, name string, width int, height int, dungeonTiles []chainTypes.DungeonTile, signature string) error {

	_ = signature

	// Off-chain checks first
	if width > chainTypes.DUNGEON_MAX_WIDTH {
		return errors.New(fmt.Sprintf("Dungeon width %d exceeds max allowed width of %d", width, chainTypes.DUNGEON_MAX_WIDTH))
	}

	if height > chainTypes.DUNGEON_MAX_HEIGHT {
		return errors.New(fmt.Sprintf("Dungeon width %d exceeds max allowed width of %d", height, chainTypes.DUNGEON_MAX_HEIGHT))
	}

	tileCount := len(dungeonTiles)

	// Allocate an array of objects to write to chain
	writeObjects := make([]chainTypes.ChainObject, 0, tileCount+2) // # of tiles + dungeon + profile

	//// On-chain checks now
	// Get the player's profile (for deducting cost)
	var playerProfile = new(chainTypes.PlayerProfile)
	var playerVault = new(chainTypes.PlayerVault)

	_, err := ReadChainObjectByParts(ctx, chainTypes.CreatePlayerProfileIdParts(owner), playerProfile)

	if err != nil {
		return err
	}

	// Get the player's vault (for validating tiles)
	_, err = ReadChainObjectByParts(ctx, chainTypes.CreatePlayerVaultIdParts(owner), playerVault)

	if err != nil {
		return err
	}

	// Ensure player has enough gold/power to create the dungeon
	if playerProfile.Gold < chainTypes.COST_CREATE_DUNGEON {
		return errors.New(fmt.Sprintf("Player has insufficient gold (%d) to create a dungeon for %d", playerProfile.Gold, chainTypes.COST_CREATE_DUNGEON))
	} else {
		playerProfile.Gold -= chainTypes.COST_CREATE_DUNGEON
	}

	// Ensure dungeon doesn't exist already
	_, ok := playerProfile.Dungeons[name]

	if ok {
		return errors.New(fmt.Sprintf("Player already has a dungeon named %s", name))
	}

	// Make the dungeon so we can populate it as we iterate through the tiles
	newDungeon := &chainTypes.Dungeon{
		Name:   name,
		Owner:  chainTypes.WalletAddress(owner),
		Power:  100,
		Height: height,
		Width:  width,
		Tiles:  make([]chainTypes.DungeonTile, 0, tileCount),
	}

	// Build a map for verifying uniqueness and reverse lookup
	dungeonTileLookup := make(map[string]*chainTypes.DungeonTile, tileCount)

	// Build a list of IDs for fetching
	tileIds := make([]string, 0, tileCount)

	// Build a map for looking up fetched tiles by ID
	tileLookup := make(map[string]*chainTypes.Tile, tileCount)

	for _, dt := range dungeonTiles {
		// Check for duplicates
		_, ok = dungeonTileLookup[dt.TileId]

		if ok {
			return errors.New(fmt.Sprintf("Duplicate tile in dungeon definition: %s", dt.TileId))
		}

		dungeonTileLookup[dt.TileId] = &dt
		tileIds = append(tileIds, dt.TileId)
	}

	// Now that we've confirmed the tiles aren't duplicated we'll
	// spend the cycles to pull them from chain and do chain-state validation

	tiles, _, err := getTiles(ctx, tileIds)

	if err != nil {
		return err
	}

	for _, t := range tiles {
		// Check that it's owned by the player
		if t.Owner != playerProfile.ExternalAddress {
			return errors.New(fmt.Sprintf("Player does not own tile being used in dungeon: %s", t.Id))
		}

		// Check that it's available
		if t.ObjectStatus != chainTypes.Available {
			return errors.New(fmt.Sprintf("Tile is not available: %s", t.Id))
		} else {
			// It is, so mark it as in use
			t.ObjectStatus = chainTypes.InUse
		}

		// Add to lookup
		tileLookup[t.Id] = t

		// Save the updates of the tile to chain and add it to the dungeon
		writeObjects = append(writeObjects, t)
		newDungeon.Tiles = append(newDungeon.Tiles, *dungeonTileLookup[t.Id])
	}

	// We performed edits to the dungeon and profile above

	// Sort dungeon tiles before saving
	// Ground, then monsters, then loot
	sort.Slice(newDungeon.Tiles, func(l, r int) bool {
		leftTile := tileLookup[newDungeon.Tiles[l].TileId]
		rightTile := tileLookup[newDungeon.Tiles[r].TileId]

		if newDungeon.Tiles[l].Y == newDungeon.Tiles[r].Y {
			if newDungeon.Tiles[l].X == newDungeon.Tiles[r].X {
				return leftTile.TileType < rightTile.TileType
			} else {
				return newDungeon.Tiles[l].X < newDungeon.Tiles[r].X
			}
		} else {
			return newDungeon.Tiles[l].Y < newDungeon.Tiles[r].Y
		}
	})

	writeObjects = append(writeObjects, newDungeon)

	playerProfile.Dungeons[newDungeon.Name] = true

	writeObjects = append(writeObjects, playerProfile)

	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SmartContract) ListDungeon(ctx contractapi.TransactionContextInterface,
	dungeonOwner string, dungeonName string, signature string) (*chainTypes.DungeonListing, error) {

	// Check signature
	_ = signature

	// Ensure dungeon listing does not exist
	dungeonListing := chainTypes.DungeonListing{
		Name:   dungeonName,
		Owner:  chainTypes.WalletAddress(dungeonOwner),
		Status: chainTypes.Available}

	_, err := ReadChainObjectByRef(ctx, &dungeonListing)

	if err == nil {
		// This means the dungeon already exists
		return nil, errors.New(fmt.Sprintf(
			"Dungeon %s:%s is already listed",
			dungeonOwner, dungeonName,
		))
	}

	// Ensure dungeon exists
	dungeon := new(chainTypes.Dungeon)
	_, err = ReadChainObjectByParts(ctx, chainTypes.CreateDungeonIdParts(dungeonOwner, dungeonName), dungeon)

	if err != nil {
		// This mean the dungeon doesn't exist
		return nil, err
	}

	// Ensure dungeon has enough power to start
	if dungeon.Power < chainTypes.DUNGEON_START_POWER {
		return nil, errors.New(fmt.Sprintf(
			"Dungeon %s:%s only has %f power but %f is required to list",
			dungeon.Owner, dungeon.Name, dungeon.Power, chainTypes.DUNGEON_START_POWER))
	}

	// Write dungeon listing
	WriteChainObject(ctx, dungeonListing)

	return &dungeonListing, nil

}

func (s *SmartContract) StartDungeon(ctx contractapi.TransactionContextInterface,
	dungeonOwner string, dungeonName string,
	characterName string, player string,
	signature string) (int64, error) {

	// Check signature
	_ = signature

	// Make sure the dungeon exists
	dungeon := &chainTypes.Dungeon{
		Owner: chainTypes.WalletAddress(dungeonOwner),
		Name:  dungeonName,
	}

	_, err := ReadChainObjectByRef(ctx, dungeon)

	if err != nil {
		return 0, err
	}

	// Make sure it's listed
	dungeonListing := &chainTypes.DungeonListing{
		Name:  dungeonName,
		Owner: chainTypes.WalletAddress(dungeonOwner),
	}

	_, err = ReadChainObjectByRef(ctx, dungeonListing)

	if err != nil {
		return 0, err
	}

	// Make sure it's available
	if dungeonListing.Status != chainTypes.Available {
		return 0, errors.New(fmt.Sprintf("Dungeon %s:%s is not available", dungeonOwner, dungeonName))
	} else {
		dungeonListing.Status = chainTypes.InUse
	}

	// Get Character
	character, err := s.GetCharacterByName(ctx, characterName)

	if err != nil {
		return 0, err
	}

	// Get profile to ensure character ownership & pay for match
	profile := chainTypes.PlayerProfile{ExternalAddress: chainTypes.WalletAddress(player)}
	_, err = ReadChainObjectByRef(ctx, &profile)

	if err != nil {
		return 0, err
	}

	// Make sure the caller owns the character
	// TODO: Actually check this in the signature
	if character.Owner != chainTypes.WalletAddress(player) {
		return 0, fmt.Errorf("Character %q is owned by %q and not by caller %q", characterName, character.Owner, player)
	}

	// Make sure the user can pay the fee
	if profile.Gold < chainTypes.COST_START_DUNGEON {
		return 0, errors.New(fmt.Sprintf(
			"Player does not have enough gold to start the dungeon: %d but %d required",
			profile.Gold, chainTypes.COST_START_DUNGEON))

	} else {
		// Charge the gold for starting the dungeon
		profile.Gold -= chainTypes.COST_START_DUNGEON
	}

	timestamp, err := ctx.GetStub().GetTxTimestamp()
	_ = timestamp

	// Create the match
	dungeonMatch := chainTypes.DungeonMatch{
		DungeonOwner:        dungeonOwner,
		DungeonName:         dungeonName,
		CharacterName:       characterName,
		Player:              player,
		DungeonMultiplier:   1.0,
		CharacterMultiplier: 1.0,
		StartingPower:       character.Power,
		// StartTime:           int64(timestamp.GetSeconds()), // TODO: Remove this so we have a real start time!
	}

	// Write the updated profile and dungeonMatch
	writeObjects := make([]chainTypes.ChainObject, 0, 2)

	writeObjects = append(writeObjects, profile)

	writeObjects = append(writeObjects, dungeonMatch)

	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return 0, err
		}
	}

	return dungeonMatch.StartTime, nil
}

func (s *SmartContract) ScoreDungeon(ctx contractapi.TransactionContextInterface,
	dungeonOwner string, dungeonName string, startTime int64,
	moves []int, signature string) (*chainTypes.DungeonMatch, error) {

	_ = signature

	// Make sure the dungeon listing exists
	dungeonListing := chainTypes.DungeonListing{
		Owner: chainTypes.WalletAddress(dungeonOwner),
		Name:  dungeonName,
	}

	// Get the dungeon listing
	_, err := ReadChainObjectByRef(ctx, &dungeonListing)

	if err != nil {
		return nil, err
	}

	dungeon := chainTypes.Dungeon{
		Owner: chainTypes.WalletAddress(dungeonOwner),
		Name:  dungeonName,
	}

	_, err = ReadChainObjectByRef(ctx, &dungeon)

	if err != nil {
		return nil, err
	}

	// Get the dungeon match
	dungeonMatch := chainTypes.DungeonMatch{
		DungeonOwner: dungeonOwner,
		DungeonName:  dungeonName,
		StartTime:    startTime,
	}

	_, err = ReadChainObjectByRef(ctx, &dungeonMatch)

	if err != nil {
		return nil, err
	}

	// Get the character
	character, err := s.GetCharacterByName(ctx, dungeonMatch.CharacterName)

	if err != nil {
		return nil, err
	}

	// Make sure the caller is the player running the dungeon
	// TODO: Improve this to check the actual identity of the caller
	if character.Owner != chainTypes.WalletAddress(signature) {
		return nil, errors.New(fmt.Sprintf(
			"Character %s is not owner by called %s",
			character.Name, signature,
		))
	}

	// Get playerProfile
	playerProfile := chainTypes.PlayerProfile{
		ExternalAddress: character.Owner,
	}

	_, err = ReadChainObjectByRef(ctx, &playerProfile)

	if err != nil {
		return nil, err
	}

	tileIds := make([]string, 0, len(dungeon.Tiles))

	for _, dt := range dungeon.Tiles {
		tileIds = append(tileIds, dt.TileId)
	}

	tiles, _, err := getTiles(ctx, tileIds)
	if err != nil {
		return nil, err
	}

	// Convert tiles to map

	tilesMap := make(map[string]*chainTypes.Tile, len(tiles))

	for _, t := range tiles {
		tilesMap[t.Id] = t
	}

	timestamp, err := ctx.GetStub().GetTxTimestamp()

	if err != nil {
		return nil, err
	}

	dungeonMatch.EndTime = timestamp.Seconds

	// Validate moves
	err = dungeonMatch.Simulate(&dungeon, tilesMap, character, moves)

	if err != nil {
		return nil, err
	}

	dungeonMatch.PlayedTiles = moves

	writeObjects := make([]chainTypes.ChainObject, 0, 4) // # Player, Dungeon, Match, Listing

	if dungeon.Power < 50 {
		// Character accumulates gold, player collects
		character.Gold += dungeonMatch.Reward
	} else {
		// Dungeon accumulates gold, player collects
		dungeon.Gold += dungeonMatch.Reward
	}

	// If dungeon still has enough power, mark it as available
	// otherwise, remove the listing
	// Ensure dungeon has enough power to start
	if dungeon.Power < chainTypes.DUNGEON_START_POWER {
		// Need to un-list the dungeon, the player needs to recharge it
		err = DeleteChainObject(ctx, dungeonListing)

		if err != nil {
			return nil, err
		}
	} else {
		dungeonListing.Status = chainTypes.Available
		writeObjects = append(writeObjects, dungeonListing)
	}

	writeObjects = append(writeObjects, dungeon)
	writeObjects = append(writeObjects, dungeonMatch)
	writeObjects = append(writeObjects, character)

	for _, o := range writeObjects {
		err = WriteChainObject(ctx, o)
		if err != nil {
			return nil, err
		}
	}

	return &dungeonMatch, nil

}

func (s *SmartContract) RechargeDungeon(ctx contractapi.TransactionContextInterface,
	dungeonOwner string, dungeonName string, signature string) error {

	_ = signature

	// Get Dungeon
	dungeon := chainTypes.Dungeon{Owner: chainTypes.WalletAddress(dungeonOwner), Name: dungeonName}

	_, err := ReadChainObjectByRef(ctx, &dungeon)

	if err != nil {
		return err
	}

	// No need to get the listing. It's either already listed and has enough power
	// or it's not listed and thus we don't need to update it

	// Get Profile
	profile := chainTypes.PlayerProfile{ExternalAddress: dungeon.Owner}

	_, err = ReadChainObjectByRef(ctx, &profile)

	if err != nil {
		return err
	}

	if dungeon.Power >= chainTypes.DUNGEON_START_POWER {
		return errors.New("Dungeon does not need to be recharged")
	}

	repairCost := int(math.Ceil(float64(chainTypes.DUNGEON_START_POWER - dungeon.Power)))

	if profile.Gold < repairCost {
		return errors.New(fmt.Sprintf("Player does not have enough gold (%d) to recharge the dungeon (%d)", profile.Gold, repairCost))
	}

	profile.Gold -= repairCost
	dungeon.Power = chainTypes.DUNGEON_START_POWER

	err = WriteChainObject(ctx, profile)

	if err != nil {
		return err
	}

	err = WriteChainObject(ctx, dungeon)

	if err != nil {
		return err
	}

	return nil

}

func (s *SmartContract) RechargeCharacter(ctx contractapi.TransactionContextInterface,
	characterOwner string, characterName string, signature string) error {

	_ = signature

	// TODO: owner should not be necessary. It should be pulled from caller / signature

	// Get Character
	character := chainTypes.Character{Name: characterName}

	_, err := ReadChainObjectByRef(ctx, &character)

	if err != nil {
		return err
	}

	if characterOwner != string(character.Owner) {
		return errors.New("Character is not owned by the caller")
	}

	// Get Profile
	profile := chainTypes.PlayerProfile{ExternalAddress: character.Owner}

	_, err = ReadChainObjectByRef(ctx, &profile)

	if err != nil {
		return err
	}

	if character.Power >= chainTypes.DUNGEON_START_POWER {
		return errors.New("Character does not need to be recharged")
	}

	repairCost := chainTypes.DUNGEON_START_POWER - int(character.Power)

	if profile.Gold < repairCost {
		return errors.New(fmt.Sprintf("Player does not have enough gold (%d) to recharge the dungeon (%d)", profile.Gold, repairCost))
	}

	profile.Gold -= repairCost
	character.Power = chainTypes.DUNGEON_START_POWER

	err = WriteChainObject(ctx, profile)

	if err != nil {
		return err
	}

	err = WriteChainObject(ctx, character)

	if err != nil {
		return err
	}

	return nil

}

func (s *SmartContract) WipeData(ctx contractapi.TransactionContextInterface) error {
	ctxStub := ctx.GetStub()
	for _, key := range chainTypes.GetAllChainKeys() {
		// Get an iterator for the key
		iterator, err := ctxStub.GetStateByPartialCompositeKey(key, []string{})

		if err != nil {
			return err
		}

		for iterator.HasNext() {
			result, err := iterator.Next()

			if err != nil {
				return err
			}

			ctxStub.DelState(result.Key)
		}
	}

	return nil

}

func (s *SmartContract) NetworkVersion() int {
	return 2
}
