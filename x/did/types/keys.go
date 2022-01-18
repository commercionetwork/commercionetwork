package types

import "time"

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
	MemStoreKey = "mem_did"

	// Version defines the current version the IBC module supports
	Version = "did-1"

	// PortID is the default port id that module binds to
	PortID = "did"

	// Identity prefix
	IdentitiesStorePrefix = StoreKey + ":identities:"

	// Context
	ContextDidV1 = "https://www.w3.org/ns/did/v1"

	MsgTypeSetDid = "MsgSetDid"

	// --------------
	// --- Queries
	// --------------

	QueryResolveDid = "identities"

	// --------------
	// --- KeyTypes required for the Documents module
	// --------------
	MultibaseCodeBase64 = 'm'

	RsaVerificationKey2018 = "RsaVerificationKey2018"
	RsaSignature2018       = "RsaSignature2018"

	RsaVerificationKey2018NameSuffix = "#keys-1"
	RsaSignature2018NameSuffix       = "#keys-2"

	// Lenght Limits
	serviceLenghtLimitID              = 56
	serviceLenghtLimitType            = 56
	serviceLenghtLimitServiceEndpoint = 256

	// XML Datetime normalized to UTC 00:00:00 and without sub-second decimal precision
	ComplaintW3CTime = time.RFC3339
)

var (

	// https://www.w3.org/TR/did-spec-registries/#verification-method-types
	verificationMethodTypes = []string{
		"Ed25519Signature2018",
		"Ed25519VerificationKey2018",
		RsaSignature2018,       // documents
		RsaVerificationKey2018, // documents
		"EcdsaSecp256k1Signature2019",
		"EcdsaSecp256k1VerificationKey2019",
		"EcdsaSecp256k1RecoverySignature2020",
		"EcdsaSecp256k1RecoveryMethod2020",
		"JsonWebSignature2020",
		"JwsVerificationKey2020",
		"GpgSignature2020",
		"GpgVerificationKey2020",
		"JcsEd25519Signature2020",
		"JcsEd25519Key2020",
		"BbsBlsSignature2020",      // vca
		"BbsBlsSignatureProof2020", // vca
		"Bls12381G1Key2020",
		"Bls12381G2Key2020",
	}
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("did-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
