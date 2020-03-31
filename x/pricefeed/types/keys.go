package types

const (
	ModuleName   = "pricefeed"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetPrice       = "setPrice"
	MsgTypeAddOracle      = "addOracle"
	MsgTypeBlacklistDenom = "blacklistDenom"

	QueryGetCurrentPrice      = "price"
	QueryGetCurrentPrices     = "prices"
	QueryGetOracles           = "oracles"
	QueryGetBlacklistedDenoms = "blacklistedDenom"

	AssetsStoreKey = StoreKey + ":assets:"

	CurrentPricesPrefix = StoreKey + ":currentPrices:"
	RawPricesPrefix     = StoreKey + ":rawPrices:"

	OraclePrefix = StoreKey + ":oracles:"

	DenomBlacklistKey = StoreKey + ":denomBlacklist:"
)
