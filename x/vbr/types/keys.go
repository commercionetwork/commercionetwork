package types

const (
	ModuleName   = "vbr"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	QueryBlockRewardsPoolFunds = "funds"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"

	PoolStoreKey  = StoreKey + ":pool:"
	RewardRateKey = StoreKey + ":rewardRate"
)
