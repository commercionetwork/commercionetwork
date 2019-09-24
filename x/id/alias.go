package id

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
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

	NewMsgSetIdentity = types.NewMsgSetIdentity

	TestSetup          = keeper.SetupTestInput
	TestDidDocumentUri = keeper.TestDidDocumentUri
	TestOwnerAddress   = keeper.TestOwnerAddress
)

type (
	Keeper = keeper.Keeper

	Identity = types.Identity

	MsgSetIdentity = types.MsgSetIdentity
)
