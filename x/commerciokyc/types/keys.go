package types

const (
	ModuleName   = "commerciokyc"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_" + ModuleName

	// --- Keeper
	MembershipsStorageKey  = StoreKey + ":membership:"
	TrustedSignersStoreKey = StoreKey + ":signers"
	InviteStorePrefix      = "invite:"

	// --- Messages
	MsgTypeInviteUser                = "inviteUser"
	MsgTypesDepositIntoLiquidityPool = "depositIntoLiquidityPool"
	MsgTypeAddTsp                    = "addTsp"
	MsgTypeRemoveTsp                 = "removeTsp"
	MsgTypeBuyMembership             = "buyMembership"
	MsgTypeSetMembership             = "setMembership"
	MsgTypeRemoveMembership          = "removeMembership"

	QueryGetInvites                 = "invites"
	QueryGetInvite                  = "invite"
	QueryGetTrustedServiceProviders = "tsps"
	QueryGetPoolFunds               = "poolFunds"
	QueryGetMembership              = "membership"
	QueryGetMemberships             = "memberships"
	QueryGetTspMemberships          = "sold"
)
