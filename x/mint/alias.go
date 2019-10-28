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
	NewHandler    = keeper.NewHandler
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgOpenCdp  = types.NewMsgOpenCdp
	NewMsgCloseCdp = types.NewMsgCloseCdp
)

type (
	Keeper = keeper.Keeper

	Cdp        = types.Cdp
	CdpRequest = types.CdpRequest

	MsgOpenCdp  = types.MsgOpenCdp
	MsgCloseCdp = types.MsgCloseCdp
)
