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

	TestSetup        = keeper.SetupTestInput
	TestDidDocument  = keeper.TestDidDocument
	TestOwnerAddress = keeper.TestOwnerAddress

	ValidateStatus = types.ValidateStatus
)

type (
	Keeper = keeper.Keeper

	Identity                = types.Identity
	DidDepositRequest       = types.DidDepositRequest
	DidDepositRequestStatus = types.DidDepositRequestStatus
	DidPowerUpRequest       = types.DidPowerUpRequest
	DidPowerUpRequestStatus = types.DidPowerUpRequestStatus

	MsgSetIdentity                   = types.MsgSetIdentity
	MsgRequestDidDeposit             = types.MsgRequestDidDeposit
	MsgChangeDidDepositRequestStatus = types.MsgChangeDidDepositRequestStatus
	MsgRequestDidPowerUp             = types.MsgRequestDidPowerUp
	MsgChangeDidPowerUpRequestStatus = types.MsgChangeDidPowerUpRequestStatus
)
