package commerciodocs

import (
	"commercio-network/x/commerciodocs/internal/keeper"
	"commercio-network/x/commerciodocs/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper      = keeper.NewKeeper
	NewQuerier     = keeper.NewQuerier
	RegisterCodec  = types.RegisterCodec
	NewMsgStoreDoc = types.NewMsgStoreDocument
	NewMsgShareDoc = types.NewMsgShareDocument

	ModuleCdc = types.ModuleCdc
)

type (
	Keeper      = keeper.Keeper
	MsgStoreDoc = types.MsgStoreDocument
	MsgShareDoc = types.MsgShareDocument
)
