package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolStoreKey = ":liquidityPool:"
	CDPStoreKey           = ":cdp:"
	CreditsDenomStoreKey  = "uccc"

	QueryGetCDP  = "CDP"
	QueryGetCDPs = "CDPs"

	MsgTypeOpenCDP  = "openCDP"
	MsgTypeCloseCDP = "closeCDP"
)
