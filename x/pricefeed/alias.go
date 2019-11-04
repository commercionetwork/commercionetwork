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
	NewHandler    = keeper.NewHandler
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewCurrentPrice = types.NewPrice

	NewMsgSetPrice  = types.NewMsgSetPrice
	NewMsgAddOracle = types.NewMsgAddOracle
)

type (
	Keeper = keeper.Keeper

	Price  = types.Price
	Prices = types.Prices

	MsgSetPrice  = types.MsgSetPrice
	MsgAddOracle = types.MsgAddOracle
)
