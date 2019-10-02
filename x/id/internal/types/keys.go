package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	QueryResolveDid            = "identities"
	QueryResolveDepositRequest = "depositRequest"
	QueryResolvePowerUpRequest = "powerUpRequest"

	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusCanceled = "canceled"

	MsgTypeSetIdentity                   = "setIdentity"
	MsgTypeRequestDidDeposit             = "requestDidDeposit"
	MsgTypeChangeDidDepositRequestStatus = "editDidDepositRequest"
	MsgTypeRequestDidPowerUp             = "requestDidPowerUp"
	MsgTypeChangeDidPowerUpRequestStatus = "editDidPowerUpRequest"

	IdentitiesStorePrefix        = StoreKey + ":identities:"
	DidDepositRequestStorePrefix = StoreKey + "depositRequest"
	DidPowerUpRequestStorePrefix = StoreKey + "PowerUpRequest"
)
