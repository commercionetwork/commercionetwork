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
	StatusApproved = types.StatusApproved
	StatusCanceled = types.StatusCanceled
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

	DidDocument         = types.DidDocument
	Proof               = types.Proof
	PubKey              = types.PubKey
	PubKeys             = types.PubKeys
	DidPowerUpRequest   = types.DidPowerUpRequest
	RequestStatus       = types.RequestStatus
	PowerUpRequestProof = types.PowerUpRequestProof

	// ---------------
	// --- Messages
	// ---------------

	MsgSetIdentity         = types.MsgSetIdentity
	MsgRequestDidPowerUp   = types.MsgRequestDidPowerUp
	MsgChangePowerUpStatus = types.MsgChangePowerUpStatus
)
