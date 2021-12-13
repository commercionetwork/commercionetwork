package types

const (
	ModuleName = "vbr"
	StoreKey = ModuleName
	RouterKey = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey = "mem_vbr"

	// this line is used by starport scaffolding # ibc/keys/name
	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"
	MsgTypeSetRewardRate             = "setRewardRate"
	MsgTypeSetAutomaticWithdraw      = "setAutomaticWithdraw"
	MsgTypeSetVbrParams              = "setVbrParams"
	PoolStoreKey                     = StoreKey + ":pool:"
	RewardRateKey                    = StoreKey + ":rewardRate"
	AutomaticWithdraw                = StoreKey + ":automaticWithdraw"

	QueryBlockRewardsPoolFunds = "funds"
	QueryRewardRate            = "reward_rate"
	QueryAutomaticWithdraw     = "automatic_withdraw"
	QueryVbrParams			   = "params"

	EpochWeek	= "week"
	EpochDay	= "day"
	EpochMinute = "minute"
)

// this line is used by starport scaffolding # ibc/keys/port

func KeyPrefix(p string) []byte {
	return []byte(p)
}
