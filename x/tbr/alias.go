package tbr

import (
	"github.com/commercionetwork/commercionetwork/x/tbr/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier

	TestSetup = keeper.SetupTestInput

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper                        = keeper.Keeper
	MsgIncrementsBlockRewardsPool = types.MsgIncrementBlockRewardsPool
)
