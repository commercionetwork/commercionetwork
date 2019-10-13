package accreditations

import (
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	// --- Messages
	NewMsgInviteUser      = types.NewMsgInviteUser
	NewMsgSetUserVerified = types.NewMsgSetUserVerified

	// --- Tests
	GetTestInput     = keeper.GetTestInput
	TestUser         = keeper.TestUser
	TestInviteSender = keeper.TestInviteSender
	TestTsp          = keeper.TestTsp
	TestTimestamp    = keeper.TestTimestamp
)

type (
	Keeper = keeper.Keeper

	Invite      = types.Invite
	Credential  = types.Credential
	Credentials = types.Credentials

	MsgInviteUser               = types.MsgInviteUser
	MsgSetUserVerified          = types.MsgSetUserVerified
	MsgDepositIntoLiquidityPool = types.MsgDepositIntoLiquidityPool
	MsgAddTrustedSigner         = types.MsgAddTsp
)
