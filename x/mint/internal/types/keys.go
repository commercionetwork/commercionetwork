package types

const (
	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	CdpStorePrefix       = ":cdp:"
	CreditsDenomStoreKey = "creditsDenom"

	QueryGetCdp  = "cdp"
	QueryGetCdps = "cdps"

	MsgTypeOpenCdp  = "openCdp"
	MsgTypeCloseCdp = "closeCdp"
)
