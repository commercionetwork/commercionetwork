package types

const (
	ModuleName   = "accreditations"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolKey      = StoreKey + ":liquidityPool:"
	TrustworthySignersKey = StoreKey + ":signers:"

	MsgTypeSetAccrediter             = "setAccrediter"
	MsgTypeDistributeReward          = "distributeReward"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"

	QueryGetAccrediter = "accrediter"
	QueryGetSigners    = "signers"
	QueryGetPoolFunds  = "getPoolFunds"
)
