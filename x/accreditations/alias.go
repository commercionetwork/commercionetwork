package accreditations

import (
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
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
	Keeper = keeper.Keeper

	MsgSetAccrediter            = types.MsgSetAccrediter
	MsgDistributeReward         = types.MsgDistributeReward
	MsgDepositIntoLiquidityPool = types.MsgDepositIntoLiquidityPool
)
