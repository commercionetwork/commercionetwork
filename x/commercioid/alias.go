package commercioid

import (
	"commercio-network/x/commercioid/internal/keeper"
	"commercio-network/x/commercioid/internal/types"
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
