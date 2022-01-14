package types

const (
	ModuleName = "vbr"
	StoreKey = ModuleName
	RouterKey = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey = "mem_vbr"

	// this line is used by starport scaffolding # ibc/keys/name
	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"
	MsgTypeSetParams              = "setParams"
	PoolStoreKey                     = StoreKey + ":pool:"

	QueryBlockRewardsPoolFunds = "funds"
	QueryParams			   = "params"

	EpochWeek	= "week"
	EpochDay	= "day"
	EpochMinute = "minute"
	EpochMonthly = "monthly"
)

// this line is used by starport scaffolding # ibc/keys/port

func KeyPrefix(p string) []byte {
	return []byte(p)
}
