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

	MsgTypeSetIdentity           = "setIdentity"
	MsgTypeRequestDidDeposit     = "requestDidDeposit"
	MsgTypeEditDidDepositRequest = "editDidDepositRequest"
	MsgTypeRequestDidPowerUp     = "requestDidPowerUp"
	MsgTypeEditDidPowerUpRequest = "editDidPowerUpRequest"

	IdentitiesStorePrefix        = StoreKey + ":identities:"
	DidDepositRequestStorePrefix = StoreKey + "depositRequest"
	DidPowerUpRequestStorePrefix = StoreKey + "PowerUpRequest"
)
