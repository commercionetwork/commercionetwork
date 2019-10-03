package id

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	StatusApproved = types.StatusApproved
	StatusRejected = types.StatusRejected
	StatusCanceled = types.StatusCanceled
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgSetIdentity = types.NewMsgSetIdentity

	NewMsgRequestDidDeposit           = types.NewMsgRequestDidDeposit
	NewMsgInvalidateDidDepositRequest = types.NewMsgInvalidateDidDepositRequest

	NewMsgRequestDidPowerUp           = types.NewMsgRequestDidPowerUp
	NewMsgInvalidateDidPowerUpRequest = types.NewMsgInvalidateDidPowerUpRequest

	// Test
	TestSetup             = keeper.SetupTestInput
	TestDidDocument       = keeper.TestDidDocument
	TestOwnerAddress      = keeper.TestOwnerAddress
	TestDidDepositRequest = keeper.TestDidDepositRequest
	TestDidPowerUpRequest = keeper.TestDidPowerUpRequest
)

type (
	Keeper = keeper.Keeper

	Identity          = types.Identity
	DidDepositRequest = types.DidDepositRequest
	DidPowerUpRequest = types.DidPowerUpRequest
	RequestStatus     = types.RequestStatus

	// ---------------
	// --- Messages
	// ---------------

	MsgSetIdentity                 = types.MsgSetIdentity
	MsgRequestDidDeposit           = types.MsgRequestDidDeposit
	MsgInvalidateDidDepositRequest = types.MsgInvalidateDidDepositRequest
	MsgRequestDidPowerUp           = types.MsgRequestDidPowerUp
	MsgInvalidateDidPowerUpRequest = types.MsgInvalidateDidPowerUpRequest
	MsgWithdrawDeposit             = types.MsgWithdrawDeposit
	MsgPowerUpDid                  = types.MsgPowerUpDid
)
