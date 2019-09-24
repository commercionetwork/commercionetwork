package types

const (
	ModuleName   = "accreditations"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	LiquidityPoolKey       = StoreKey + ":liquidityPool:"
	TrustedSignersStoreKey = StoreKey + ":signers:"

	MsgTypeSetAccrediter             = "setAccrediter"
	MsgTypeDistributeReward          = "distributeReward"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"
	MsgTypeAddTrustedSigner          = "addTrustedSigner"

	QueryGetAccrediter = "accrediter"
	QueryGetSigners    = "signers"
	QueryGetPoolFunds  = "getPoolFunds"
)
