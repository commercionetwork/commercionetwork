package memberships

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	//function aliases
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	TestSetup     = keeper.SetupTestInput
	RegisterCodec = types.RegisterCodec

	//variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper = keeper.Keeper

	Minters = types.Minters

	MsgAssignMembership = types.MsgAssignMembership
)
