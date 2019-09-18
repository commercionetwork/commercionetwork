package mint

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
)

const (
	ModuleName          = types.ModuleName
	StoreKey            = types.StoreKey
	QuerierRoute        = types.QuerierRoute
	DefaultCreditsDenom = types.DefaultCreditsDenom
	DefaultBondDenom    = types.DefaultBondDenom
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	TestSetup     = keeper.SetupTestInput
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
)

type (
	Keeper           = keeper.Keeper
	MsgDepositToken  = types.MsgDepositToken
	MsgWithdrawToken = types.MsgWithdrawToken
)
