package types

const (
	ModuleName   = "tbr"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	QueryBlockRewardsPoolFunds = "funds"

	MsgTypeIncrementBlockRewardsPool = "incrementBlockRewardsPool"

	PoolStoreKey       = StoreKey + ":pool:"
	YearlyPoolStoreKey = StoreKey + ":yearly_pool:"
	YearlyPoolRemains  = StoreKey + ":yearly_pool_remains:"
	YearNumberStoreKey = StoreKey + ":year_number:"
)
