package types

const (
	ModuleName   = "commerciomint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	CdpStorePrefix       = ":cdp:"
	CreditsDenomStoreKey = "CreditsDenom"
	CollateralRateKey    = StoreKey + ":collateralRate"

	QueryGetCdp         = "cdp"
	QueryGetCdps        = "cdps"
	QueryCollateralRate = "collateral_rate"

	MsgTypeOpenCdp              = "openCdp"
	MsgTypeCloseCdp             = "closeCdp"
	MsgTypeSetCdpCollateralRate = "setCdpCollateralRate"
)
