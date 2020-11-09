package types

const (
	ModuleName   = "accreditations"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	// --- Keeper
	MembershipsStorageKey  = StoreKey + ":storage:"
	TrustedSignersStoreKey = StoreKey + ":signers"
	InviteStorePrefix      = "invite:"
	CredentialsStorePrefix = "credentials:"

	// --- Messages
	MsgTypeInviteUser                = "inviteUser"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"
	MsgTypeAddTsp                    = "addTsp"
	MsgTypeRemoveTsp                 = "removeTsp"
	MsgTypeBuyMembership             = "buyMembership"
	MsgTypeSetMembership             = "setMembership"
	MsgTypeRemoveMembership          = "removeMembership"

	QueryGetInvites                 = "invites"
	QueryGetTrustedServiceProviders = "tsps"
	QueryGetPoolFunds               = "poolFunds"
	QueryGetMembership              = "membership"
	QueryGetMemberships             = "memberships"
	QueryGetTspMemberships          = "sold"
)
