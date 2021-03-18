package types

const (
	ModuleName   = "commerciomint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	EtpStorePrefix       = StoreKey + ":etp:"
	CreditsDenomStoreKey = "CreditsDenom"
	CollateralRateKey    = StoreKey + ":collateralRate"
	FreezePeriodKey      = StoreKey + ":freezePeriod"

	QueryGetEtps        = "etps"
	QueryConversionRate = "conversion_rate"
	QueryFreezePeriod   = "freeze_period"

	MsgTypeMintCCC              = "mintCCC"
	MsgTypeBurnCCC              = "burnCCC"
	MsgTypeSetCCCConversionRate = "setEtpsConversionRate"
	MsgTypeSetCCCFreezePeriod   = "setEtpsFreezePeriod"
)
