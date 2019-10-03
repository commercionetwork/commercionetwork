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

	NewMsgSetIdentity                   = types.NewMsgSetIdentity
	NewMsgRequestDidDeposit             = types.NewMsgRequestDidDeposit
	NewMsgChangeDidDepositRequestStatus = types.NewMsgChangeDidDepositRequestStatus
	NewMsgRequestDidPowerUp             = types.NewMsgRequestDidPowerUp
	NewMsgChangeDidPowerUpRequestStatus = types.NewMsgChangeDidPowerUpRequestStatus
	ValidateStatus                      = types.ValidateStatus

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

	MsgSetIdentity                   = types.MsgSetIdentity
	MsgRequestDidDeposit             = types.MsgRequestDidDeposit
	MsgChangeDidDepositRequestStatus = types.MsgChangeDidDepositRequestStatus
	MsgRequestDidPowerUp             = types.MsgRequestDidPowerUp
	MsgChangeDidPowerUpRequestStatus = types.MsgChangeDidPowerUpRequestStatus
)
