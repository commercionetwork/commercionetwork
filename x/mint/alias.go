package mint

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	TestSetup     = keeper.SetupTestInput
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
)
