package types

const (
	ModuleName   = "government"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	GovernmentStoreKey = StoreKey + "government"
	TumblerStoreKey    = StoreKey + "tumbler"

	QueryGovernmentAddress = "governmentAddress"
	QueryTumblerAddress    = "tumblerAddress"
)
