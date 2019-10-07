package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolStorePrefix = ":liquidityPool:"
	UserCdpsStorePrefix      = ":cdp:"
	CreditsDenomStoreKey     = StoreKey + ":creditsDenom:"

	QueryGetCdp  = "Cdp"
	QueryGetCdps = "Cdps"

	MsgTypeOpenCdp  = "openCdp"
	MsgTypeCloseCdp = "closeCdp"
)
