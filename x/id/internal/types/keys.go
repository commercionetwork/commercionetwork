package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetIdentity = "setIdentity"

	IdentitiesStorePrefix = StoreKey + ":identities:"

	KeyTypeRsa       = "RsaVerificationKey2018"
	KeyTypeSecp256k1 = "Secp256k1VerificationKey2018"
	KeyTypeEd25519   = "Ed25519VerificationKey2018"
)
