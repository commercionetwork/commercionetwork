package txreward

import (
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
)

const (
	ModuleName             = types.ModuleName
	StoreKey               = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	BlockRewardsPoolPrefix = types.BlockRewardsPoolPrefix
)

var (
	//function aliases
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec

	//variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper                        = keeper.Keeper
	MsgIncrementsBlockRewardsPool = types.MsgIncrementBlockRewardsPool
)
