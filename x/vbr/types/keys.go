package types

const (
	ModuleName   = "vbr"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName

	QueryBlockRewardsPoolFunds = "funds"
	QueryRewardRate            = "reward_rate"
	QueryAutomaticWithdraw     = "automatic_withdraw"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"
	MsgTypeSetRewardRate             = "setRewardRate"
	MsgTypeSetAutomaticWithdraw      = "setAutomaticWithdraw"
	PoolStoreKey                     = StoreKey + ":pool:"
	RewardRateKey                    = StoreKey + ":rewardRate"
	AutomaticWithdraw                = StoreKey + ":automaticWithdraw"
)
