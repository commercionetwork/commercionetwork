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

	// --- Messages
	MsgTypeInviteUser = "inviteUser"

	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"
	MsgTypeAddTsp                    = "addTsp"
	MsgTypeBuyMembership             = "buyMembership"
	MsgTypeSetMembership             = "setMembership"

	QueryGetInvites                 = "invites"
	QueryGetTrustedServiceProviders = "tsps"
	QueryGetPoolFunds               = "poolFunds"
	QueryGetMembership              = "memberships"
)
