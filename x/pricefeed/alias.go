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

	NewCurrentPrice = types.NewCurrentPrice

	NewMsgSetPrice  = types.NewMsgSetPrice
	NewMsgAddOracle = types.NewMsgAddOracle

	// Testing
	TestInput      = keeper.SetupTestInput
	TestOracle1    = keeper.TestOracle1
	TestGovernment = keeper.TestGovernment
	TestRawPrice   = keeper.TestRawPrice1
)

type (
	Keeper = keeper.Keeper

	CurrentPrice = types.CurrentPrice
	RawPrice     = types.RawPrice

	MsgSetPrice  = types.MsgSetPrice
	MsgAddOracle = types.MsgAddOracle
)
