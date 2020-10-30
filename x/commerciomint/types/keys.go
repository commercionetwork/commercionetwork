package types

const (
	ModuleName   = "commerciomint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	EtpStorePrefix       = ":etp:"
	CreditsDenomStoreKey = "CreditsDenom"
	CollateralRateKey    = StoreKey + ":collateralRate"

	QueryGetEtps        = "etps"
	QueryConversionRate = "conversion_rate"

	MsgTypeMintCCC              = "mintCCC"
	MsgTypeBurnCCC              = "burnCCC"
	MsgTypeSetCCCConversionRate = "setEtpsConversionRate"
)
