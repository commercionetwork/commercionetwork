package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolPrefix = ":liquidityPool:"
	CDPSPrefix          = ":cdp:"

	CreditsDenom = "uccc"

	QueryGetCDP  = "CDP"
	QueryGetCDPs = "CDPs"

	MsgTypeOpenCDP  = "openCDP"
	MsgTypeCloseCDP = "closeCDP"
)
