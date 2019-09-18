package types

import "github.com/commercionetwork/commercionetwork/app"

const (
	DefaultBondDenom    = app.DefaultBondDenom
	DefaultCreditsDenom = "ucommerciocredits"

	ModuleName   = "mint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	MsgTypeDepositToken  = "depositToken"
	MsgTypeWithdrawToken = "withdrawToken"
)
