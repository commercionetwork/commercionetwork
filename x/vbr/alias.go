package vbr

import (
	"github.com/commercionetwork/commercionetwork/x/vbr/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/vbr/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper  = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper

	MsgIncrementsBlockRewardsPool = types.MsgIncrementBlockRewardsPool
)
