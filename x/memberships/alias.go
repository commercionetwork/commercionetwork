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
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
	RegisterCodec          = types.RegisterCodec
	NewMsgAssignMembership = types.NewMsgAssignMembership

	//variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper              = keeper.Keeper
	MsgAssignMembership = types.MsgAssignMembership
)
