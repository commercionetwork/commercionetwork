package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	UserCdpsStorePrefix  = ":cdp:"
	CreditsDenomStoreKey = StoreKey + ":creditsDenom:"

	QueryGetCdp  = "Cdp"
	QueryGetCdps = "Cdps"

	MsgTypeOpenCdp  = "openCdp"
	MsgTypeCloseCdp = "closeCdp"
)
