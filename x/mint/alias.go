package mint

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgOpenCDP  = types.NewMsgOpenCDP
	NewMsgCloseCDP = types.NewMsgCloseCDP

	//Testing
	TestInput = keeper.SetupTestInput
)

type (
	Keeper           = keeper.Keeper
	MsgDepositToken  = types.MsgOpenCDP
	MsgWithdrawToken = types.MsgCloseCDP
)
