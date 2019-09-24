package pricefeed

import (
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
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
)

type (
	Keeper       = keeper.Keeper
	MsgSetPrice  = types.MsgSetPrice
	MsgAddOracle = types.MsgAddOracle
)
