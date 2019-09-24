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
	TestSetup     = keeper.SetupTestInput
	RegisterCodec = types.RegisterCodec

	ModuleCdc        = types.ModuleCdc
	TestOwnerAddress = keeper.TestOwnerAddress
	TestDidDocument  = keeper.TestDidDocument
)

type (
	Keeper = keeper.Keeper

	DidDocument = types.DidDocument
	Identity    = types.Identity

	MsgSetIdentity = types.MsgSetIdentity
)
