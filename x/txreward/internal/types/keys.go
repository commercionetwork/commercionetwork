package types

const (
	ModuleName   = "txBlockReward"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"

	BlockRewardsPoolPrefix = StoreKey + ":bRewardsPool:"

	DefaultBondDenom = "ucommercio"
)
