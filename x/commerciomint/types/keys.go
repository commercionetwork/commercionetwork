package types

const (
	ModuleName   = "commerciomint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_" + ModuleName

	EtpStorePrefix = StoreKey + ":etp:"
	CreditsDenom   = "uccc"
	BondDenom      = "ucommercio"

	QueryGetEtpRest         = "etp"
	QueryGetallEtpsRest     = "etps"
	QueryGetEtpsByOwnerRest = "etpsOwner"
	QueryConversionRateRest = "conversion_rate"
	QueryFreezePeriodRest   = "freeze_period"
	QueryGetParamsRest      = "params"

	MsgTypeMintCCC   = "mintCCC"
	MsgTypeBurnCCC   = "burnCCC"
	MsgTypeSetParams = "setParams"
)
