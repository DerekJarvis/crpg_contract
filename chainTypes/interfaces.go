package chainTypes

type ChainObject interface {
	// Returns the chain key + ID components. Chain type should be the first string
	GetIdParts() []string
}

type HasOwner interface {
	GetOwner() WalletAddress
}
