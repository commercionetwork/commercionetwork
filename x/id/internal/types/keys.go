package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	// -----------------
	// --- Store keys
	// -----------------

	IdentitiesStorePrefix = StoreKey + ":identities:"

	DepositsPoolStoreKey           = StoreKey + "depositsPool"
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

	MsgTypeWithdrawDeposit = "withdrawDeposit"
	MsgTypePowerUpDid      = "powerUpDid"
)
