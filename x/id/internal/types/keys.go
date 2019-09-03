package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetIdentity = "setIdentity"

	IdentitiesStorePrefix = StoreKey + ":identities:"
)
