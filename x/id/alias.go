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

	TestSetup        = keeper.SetupTestInput
	TestDidDocument  = keeper.TestDidDocument
	TestOwnerAddress = keeper.TestOwnerAddress
)

type (
	Keeper = keeper.Keeper

	Identity = types.Identity

	MsgSetIdentity = types.MsgSetIdentity
)
