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
	//test vars
	NewMsgSetPrice  = types.NewMsgSetPrice
	NewMsgAddOracle = types.NewMsgAddOracle

	TestInput      = keeper.SetupTestInput
	TestOracle1    = keeper.TestOracle1
	TestGovernment = keeper.TestGovernment
	TestRawPrice   = keeper.TestRawPrice1
)

type (
	Keeper       = keeper.Keeper
	MsgSetPrice  = types.MsgSetPrice
	CurrentPrice = types.CurrentPrice
	RawPrice     = types.RawPrice
	MsgAddOracle = types.MsgAddOracle
)
