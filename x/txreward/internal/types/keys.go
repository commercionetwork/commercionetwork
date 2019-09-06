package types

const (
	ModuleName   = "txBlockReward"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	QueryBlockRewardsPoolFunds   = "funds"
	QueryBlockRewardsPoolFunders = "funders"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"

	BlockRewardsPoolPrefix        = StoreKey + ":pool:"
	BlockRewardsPoolFundersPrefix = StoreKey + ":funders:"

	DefaultBondDenom = "ucommercio"
)
