package types

const (
	ModuleName   = "accreditations"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	InviteStorePrefix      = "invite:"
	LiquidityPoolStoreKey  = StoreKey + ":liquidityPool:"
	TrustedSignersStoreKey = StoreKey + ":signers:"
	CredentialsStorePrefix = "credentials:"

	MsgTypeInviteUser                = "inviteUser"
	MsgTypeSetUserVerified           = "setUserVerified"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"
	MsgTypeAddTsp                    = "addTsp"

	QueryGetInvites                 = "invites"
	QueryGetTrustedServiceProviders = "tsps"
	QueryGetPoolFunds               = "poolFunds"
)
