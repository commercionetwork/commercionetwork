package id

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	StatusRejected = types.StatusRejected
)

var (
	NewKeeper     = keeper.NewKeeper
	NewHandler    = keeper.NewHandler
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
)

type (
	Keeper = keeper.Keeper

	DidDocument       = types.DidDocument
	DidDepositRequest = types.DidDepositRequest
	DidPowerUpRequest = types.DidPowerUpRequest
	RequestStatus     = types.RequestStatus

	// ---------------
	// --- Messages
	// ---------------

	MsgSetIdentity                 = types.MsgSetIdentity
	MsgRequestDidPowerUp           = types.MsgRequestDidPowerUp
	MsgInvalidateDidPowerUpRequest = types.MsgInvalidateDidPowerUpRequest
	MsgPowerUpDid                  = types.MsgPowerUpDid
)
