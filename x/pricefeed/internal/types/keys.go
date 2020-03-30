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

	AssetsStoreKey = StoreKey + ":assets:"

	CurrentPricesPrefix = StoreKey + ":currentPrices:"
	RawPricesPrefix     = StoreKey + ":rawPrices:"

	OraclePrefix = StoreKey + ":oracles:"

	DenomBlacklistKey = StoreKey + ":denomBlacklist:"
)
