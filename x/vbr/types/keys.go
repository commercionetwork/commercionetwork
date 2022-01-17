package types

const (
	ModuleName = "vbr"
	StoreKey = ModuleName
	RouterKey = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey = "mem_vbr"

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

func KeyPrefix(p string) []byte {
	return []byte(p)
}
