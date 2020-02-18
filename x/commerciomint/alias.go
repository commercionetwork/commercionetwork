package commerciomint

import (
	"github.com/commercionetwork/commercionetwork/x/commerciomint/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	NewHandler    = keeper.NewHandler
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgOpenCdp      = types.NewMsgOpenCdp
	NewMsgCloseCdp     = types.NewMsgCloseCdp
	RegisterInvariants = keeper.RegisterInvariants
)

type (
	Keeper = keeper.Keeper

	Cdp  = types.Cdp
	Cdps = types.Cdps

	MsgOpenCdp  = types.MsgOpenCdp
	MsgCloseCdp = types.MsgCloseCdp
)
