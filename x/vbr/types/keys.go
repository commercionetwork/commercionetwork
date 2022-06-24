package types

const (
	ModuleName   = "vbr"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_" + ModuleName

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"
	MsgTypeSetParams                 = "setParams"
	PoolStoreKey                     = StoreKey + ":pool:"

	QueryBlockRewardsPoolFunds = "funds"
	QueryParams                = "params"

	EpochWeek   = "week"
	EpochDay    = "day"
	EpochHour   = "hour"
	EpochMinute = "minute"
	EpochMonth  = "month"

	BondDenom = "ucommercio"
)

// TODO: use KeyPrefix function if needed
/* func KeyPrefix(p string) []byte {
	return []byte(p)
}*/
