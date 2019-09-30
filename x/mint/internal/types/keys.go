package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolStoreKey = ":liquidityPool:"
	CDPStoreKey           = ":cdp:"
	CreditsDenomStoreKey  = "uccc"

	MsgTypeDepositToken  = "depositToken"
	MsgTypeWithdrawToken = "withdrawToken"
)
