package types

const (
	ModuleName   = "vbr"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName

	QueryBlockRewardsPoolFunds = "funds"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"
	MsgTypeSetRewardRate             = "setRewardRate"
	MsgTypeSetAutomaticWithdraw      = "setAutomaticWithdraw"
	PoolStoreKey                     = StoreKey + ":pool:"
	RewardRateKey                    = StoreKey + ":rewardRate"
	AutomaticWithdraw                = StoreKey + ":automaticWithdraw"
)
