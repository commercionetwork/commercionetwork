package types

const (
	ModuleName   = "accreditations"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetAccrediter             = "setAccrediter"
	MsgTypeDistributeReward          = "distributeReward"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"

	QueryGetAccrediter = "accrediter"
	QueryGetSigners    = "signers"

	LiquidityPoolKey      = "liquidityPool"
	TrustworthySignersKey = "signers"
)
