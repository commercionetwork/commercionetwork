package types

const (
	ModuleName   = "pricefeed"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetPrice  = "setPrice"
	MsgTypeAddOracle = "addOracle"

	QueryGetTokenPrice = "price"
	QueryGetOracles    = "oracles"

	//CurrentPricePrefix store prefix for the current price of an asset
	CurrentPricePrefix = StoreKey + ":currentprice:"

	//OraclePrefix store prefix for the oracle accounts
	OraclePrefix = StoreKey + ":oracles:"
)
