package types

const (
	ModuleName   = "accreditations"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	// --- Keeper
	MembershipsStorageKey  = StoreKey + ":storage:"
	StableCreditsStoreKey  = StoreKey + ":stableCreditsDenom"
	TrustedSignersStoreKey = StoreKey + ":signers"
	InviteStorePrefix      = "invite:"
	CredentialsStorePrefix = "credentials:"

	// --- Messages
	MsgTypeInviteUser                = "inviteUser"
	MsgTypeSetUserVerified           = "setUserVerified"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"
	MsgTypeAddTsp                    = "addTsp"
	MsgTypeBuyMembership             = "buyMembership"

	QueryGetInvites                 = "invites"
	QueryGetTrustedServiceProviders = "tsps"
	QueryGetPoolFunds               = "poolFunds"
	QueryGetMembership              = "memberships"
)
