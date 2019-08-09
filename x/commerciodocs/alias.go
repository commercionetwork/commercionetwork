package commerciodocs

import (
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
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
	Keeper           = keeper.Keeper
	MsgShareDocument = types.MsgShareDocument
)
