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

	InviteStatusInvalid  = types.InviteStatusInvalid
	InviteStatusPending  = types.InviteStatusPending
	InviteStatusRewarded = types.InviteStatusRewarded
)

type (
	Keeper = keeper.Keeper

	Invite       = types.Invite
	Invites      = types.Invites
	Membership   = types.Membership
	Memberships  = types.Memberships
	InviteStatus = types.InviteStatus

	MsgInviteUser               = types.MsgInviteUser
	MsgDepositIntoLiquidityPool = types.MsgDepositIntoLiquidityPool
	MsgAddTrustedSigner         = types.MsgAddTsp
	MsgBuyMembership            = types.MsgBuyMembership
)
