package types

const (
	ModuleName   = "memberships"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	// --- Keeper
	NftDenom = "membership"

	// --- Messages
	MsgTypeBuyMembership = "buyMembership"

	// --- Queries
	QueryGetMembership = "memberships"
)
