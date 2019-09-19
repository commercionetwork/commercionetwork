package types

const (
	ModuleName   = "txBlockReward"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	QueryBlockRewardsPoolFunds = "funds"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"

	BlockRewardsPoolPrefix = StoreKey + ":pool:"

	DefaultBondDenom = "ucommercio"
)
