package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetIdentity           = "setIdentity"
	MsgTypeRequestDidDeposit     = "requestDidDeposit"
	MsgTypeEditDidDepositRequest = "editDidDepositRequest"
	MsgTypeRequestDidPowerup     = "requestDidPowerup"
	MsgTypeEditDidPowerupRequest = "editDidPowerupRequest"

	IdentitiesStorePrefix = StoreKey + ":identities:"
)
