package membership

import (
	"github.com/commercionetwork/commercionetwork/x/membership/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/membership/internal/types"
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
	RegisterCodec = types.RegisterCodec

	//variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper              = keeper.Keeper
	MsgAssignMembership = types.MsgAssignMembership
)
