package types

const (
	ModuleName   = "tbr"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	QueryBlockRewardsPoolFunds = "funds"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"

	PoolStoreKey        = StoreKey + ":pool:"
	RewardDenomStoreKey = StoreKey + ":reward_denom:"
)
