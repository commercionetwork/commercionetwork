package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolStoreKey = ":liquidityPool:"
	CDPStoreKey           = ":cdp:"
	CreditsDenomStoreKey  = "uccc"

	MsgTypeOpenCDP  = "openCDP"
	MsgTypeCloseCDP = "closeCDP"
)
