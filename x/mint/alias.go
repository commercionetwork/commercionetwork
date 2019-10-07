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

	NewMsgOpenCdp  = types.NewMsgOpenCdp
	NewMsgCloseCdp = types.NewMsgCloseCdp

	//Testing
	TestInput      = keeper.SetupTestInput
	TestCdpRequest = keeper.TestCdpRequest
	TestCdp        = keeper.TestCdp
)

type (
	Keeper           = keeper.Keeper
	MsgDepositToken  = types.MsgOpenCdp
	MsgWithdrawToken = types.MsgCloseCdp
	CdpRequest       = types.CdpRequest
	Cdp              = types.Cdp
)
