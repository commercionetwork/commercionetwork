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

	NewPrice   = types.NewPrice
	EmptyPrice = types.EmptyPrice

	NewMsgSetPrice  = types.NewMsgSetPrice
	NewMsgAddOracle = types.NewMsgAddOracle
)

type (
	Keeper = keeper.Keeper

	Price     = types.Price
	Prices    = types.Prices
	RawPrice  = types.RawPrice
	RawPrices = types.RawPrices

	MsgSetPrice  = types.MsgSetPrice
	MsgAddOracle = types.MsgAddOracle
)
