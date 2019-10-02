package types

const (
	ModuleName   = "id"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeSetIdentity              = "setIdentity"
	MsgTypeRequestDidDeposit        = "requestDidDeposit"
	MsgTypeSetDepositRequestHandled = "setDepositRequestHandled"
	MsgTypeEditDidDepositRequest    = "editDidDepositRequest"
	MsgTypeRequestDidPowerup        = "requestDidPowerup"
	MsgTypePowerupDid               = "powerupDid"
	MsgTypeEditDidPowerupRequest    = "editDidPowerupRequest"

	IdentitiesStorePrefix = StoreKey + ":identities:"
)
