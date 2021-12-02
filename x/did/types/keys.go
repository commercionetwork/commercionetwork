package types

const (
	// ModuleName defines the module name
	ModuleName = "did"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_id"

	// Version defines the current version the IBC module supports
	Version = "did-1"

	// PortID is the default port id that module binds to
	PortID = "did"

	// Identity prefix
	IdentitiesStorePrefix = StoreKey + ":identities:"

	// Context
	ContextDidV1 = "https://www.w3.org/ns/did/v1"

	// key types
	KeyTypeRsaVerification   = "RsaVerificationKey2018"
	KeyTypeRsaSignature      = "RsaSignatureKey2018"
	KeyTypeSecp256k1         = "Secp256k1VerificationKey2018"
	KeyTypeSecp256k12019     = "EcdsaSecp256k1VerificationKey2019"
	KeyTypeEd25519           = "Ed25519VerificationKey2018"
	KeyTypeBls12381G1Key2020 = "Bls12381G1Key2020"
	KeyTypeBls12381G2Key2020 = "Bls12381G2Key2020"

	DidPowerUpRequestStorePrefix               = StoreKey + "powerUpRequest"
	HandledPowerUpRequestsReferenceStorePrefix = StoreKey + "handledPowerUpRequestsReference"

	MsgTypeSetDid = "MsgSetDid"

	// --------------
	// --- Queries
	// --------------

	QueryResolveDid = "identities"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("id-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
