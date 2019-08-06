package commercioid

import (
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"
)

const (
	ModuleName          = types.ModuleName
	IdentitiesStoreKey  = types.IdentitiesStoreKey
	OwnersStoreKey      = types.OwnersStoreKey
	ConnectionsStoreKey = types.ConnectionsStoreKey
	QuerierRoute        = types.QuerierRoute
)

var (
	//function aliases
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
	RegisterCodec          = types.RegisterCodec
	NewMsgSetIdentity      = types.NewMsgSetIdentity
	NewMsgCreateConnection = types.NewMsgCreateConnection

	//variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper              = keeper.Keeper
	MsgSetIdentity      = types.MsgSetIdentity
	MsgCreateConnection = types.MsgCreateConnection
)
