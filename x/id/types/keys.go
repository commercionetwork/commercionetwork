package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	// -----------------
	// --- Store keys
	// -----------------

	IdentitiesStorePrefix = StoreKey + ":identities:"

	KeyTypeRsaVerification = "RsaVerificationKey2018"
	KeyTypeRsaSignature    = "RsaSignatureKey2018"
	KeyTypeSecp256k1       = "Secp256k1VerificationKey2018"
	KeyTypeSecp256k12019   = "EcdsaSecp256k1VerificationKey2019"
	KeyTypeEd25519         = "Ed25519VerificationKey2018"

	DidDepositRequestStorePrefix               = StoreKey + "depositRequest"
	DidPowerUpRequestStorePrefix               = StoreKey + "powerUpRequest"
	HandledPowerUpRequestsReferenceStorePrefix = StoreKey + "handledPowerUpRequestsReference"

	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusCanceled = "canceled"

	// --------------
	// --- Queries
	// --------------

	QueryResolveDid = "identities"

	QueryResolveDepositRequest     = "depositRequest"
	QueryResolvePowerUpRequest     = "powerUpRequest"
	QueryGetApprovedPowerUpRequest = "approvedPowerUpRequest"
	QueryGetRejectedPowerUpRequest = "rejectedPowerUpRequest"
	QueryGetPendingPowerUpRequest  = "pendingPowerUpRequest"

	// --------------
	// --- Messages
	// --------------

	MsgTypeSetIdentity         = "setIdentity"
	MsgTypeRequestDidPowerUp   = "requestDidPowerUp"
	MsgTypeChangePowerUpStatus = "changePowerUpStatus"

	// --------------
	// --- Context
	// --------------
	ContextDidV1 = "https://www.w3.org/ns/did/v1"
)
