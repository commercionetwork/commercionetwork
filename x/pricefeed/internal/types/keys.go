package types

const (
	ModuleName   = "pricefeed"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetPrice = "setPrice"

	QueryGetTokenPrice = "price"

	// Store prefix for the current price of an asset
	CurrentPricePrefix = StoreKey + ":currentprice:"
)
