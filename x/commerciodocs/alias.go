package commerciodocs

import (
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
)

const (
	ModuleName       = types.ModuleName
	OwnersStoreKey   = types.OwnersStoreKey
	MetadataStoreKey = types.MetadataStoreKey
	SharingStoreKey  = types.SharingStoreKey
	QuerierRoute     = types.QuerierRoute
)

var (
	NewKeeper      = keeper.NewKeeper
	NewQuerier     = keeper.NewQuerier
	RegisterCodec  = types.RegisterCodec
	NewMsgShareDoc = types.NewMsgShareDocument

	ModuleCdc = types.ModuleCdc
)

type (
	Keeper      = keeper.Keeper
	MsgShareDoc = types.MsgShareDocument
)
