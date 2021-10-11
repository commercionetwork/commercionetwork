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
	CreditsDenom         = "uccc"
	BondDenom            = "ucommercio"

	QueryGetEtps = "etps"

	MsgTypeMintCCC              = "mintCCC"
	MsgTypeBurnCCC              = "burnCCC"
	MsgTypeSetCCCConversionRate = "setEtpsConversionRate"
	MsgTypeSetCCCFreezePeriod   = "setEtpsFreezePeriod"
)
