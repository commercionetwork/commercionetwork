package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	// -----------------
	// --- Store keys
	// -----------------

	IdentitiesStorePrefix = StoreKey + ":identities:"

	KeyTypeRsa       = "RsaVerificationKey2018"
	KeyTypeSecp256k1 = "Secp256k1VerificationKey2018"
	KeyTypeEd25519   = "Ed25519VerificationKey2018"

	DidDepositRequestStorePrefix   = StoreKey + "depositRequest"
	DidPowerUpRequestStorePrefix   = StoreKey + "powerUpRequest"
	HandledPowerUpRequestsStoreKey = StoreKey + "handledPowerUpRequests"

	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusCanceled = "canceled"

	// --------------
	// --- Queries
	// --------------

	QueryResolveDid = "identities"

	QueryResolveDepositRequest = "depositRequest"
	QueryResolvePowerUpRequest = "powerUpRequest"

	// --------------
	// --- Messages
	// --------------

	MsgTypeSetIdentity = "setIdentity"

	MsgTypeRequestDidDeposit           = "requestDidDeposit"
	MsgTypeInvalidateDidDepositRequest = "invalidateDidDepositRequest"
	MsgTypeRequestDidPowerUp           = "requestDidPowerUp"
	MsgTypeInvalidateDidPowerUpRequest = "invalidateDidPowerUpRequest"

	MsgTypeMoveDeposit = "moveDeposit"
	MsgTypePowerUpDid  = "powerUpDid"
)
