package types

const (
	ModuleName   = "pricefeed"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetPrice  = "setPrice"
	MsgTypeAddOracle = "addOracle"

	QueryGetCurrentPrice  = "price"
	QueryGetCurrentPrices = "prices"
	QueryGetOracles       = "oracles"

	//CurrentPricesPrefix store prefix for the current price of an asset
	AssetsPrefix        = StoreKey + ":tokeninfo:"
	CurrentPricesPrefix = StoreKey + ":currentprices:"
	RawPricesPrefix     = StoreKey + ":rawprices:"

	//OraclePrefix store prefix for the oracle accounts
	OraclePrefix = StoreKey + ":oracles:"
)
