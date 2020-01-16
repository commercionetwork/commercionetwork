package memberships

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
	NewHandler = keeper.NewHandler

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMembership         = types.NewMembership
	IsMembershipTypeValid = types.IsMembershipTypeValid

	RegisterInvariants = keeper.RegisterInvariants
)

type (
	Keeper = keeper.Keeper

	Invite      = types.Invite
	Invites     = types.Invites
	Credential  = types.Credential
	Credentials = types.Credentials
	Membership  = types.Membership
	Memberships = types.Memberships

	MsgInviteUser               = types.MsgInviteUser
	MsgSetUserVerified          = types.MsgSetUserVerified
	MsgDepositIntoLiquidityPool = types.MsgDepositIntoLiquidityPool
	MsgAddTrustedSigner         = types.MsgAddTsp
	MsgBuyMembership            = types.MsgBuyMembership
)
