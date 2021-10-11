package types

const (
	// ModuleName defines the module name
	ModuleName = "vbr"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vbr"

	// this line is used by starport scaffolding # ibc/keys/name
	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"
	MsgTypeSetRewardRate             = "setRewardRate"
	MsgTypeSetAutomaticWithdraw      = "setAutomaticWithdraw"
	PoolStoreKey                     = StoreKey + ":pool:"
	RewardRateKey                    = StoreKey + ":rewardRate"
	AutomaticWithdraw                = StoreKey + ":automaticWithdraw"
)

// this line is used by starport scaffolding # ibc/keys/port

func KeyPrefix(p string) []byte {
	return []byte(p)
}
