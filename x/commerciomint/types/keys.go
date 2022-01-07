package types

const (
	ModuleName   = "commerciomint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_" + ModuleName

	EtpStorePrefix       = StoreKey + ":etp:"
	CreditsDenomStoreKey = "CreditsDenom"
	CollateralRateKey    = StoreKey + ":collateralRate"
	FreezePeriodKey      = StoreKey + ":freezePeriod"
	CreditsDenom         = "uccc"
	BondDenom            = "ucommercio"

	QueryGetEtpRest         = "etp"
	QueryGetallEtpsRest     = "etps"
	QueryGetEtpsByOwnerRest = "owner"
	QueryConversionRateRest = "conversion_rate"
	QueryFreezePeriodRest   = "freeze_period"

	MsgTypeMintCCC              = "mintCCC"
	MsgTypeBurnCCC              = "burnCCC"
	MsgTypeSetCCCConversionRate = "setEtpsConversionRate"
	MsgTypeSetCCCFreezePeriod   = "setEtpsFreezePeriod"
)
