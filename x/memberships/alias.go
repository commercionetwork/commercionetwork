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
)

type (
	Keeper = keeper.Keeper

	Invite      = types.Invite
	Credential  = types.Credential
	Credentials = types.Credentials
	Membership  = types.Membership

	MsgInviteUser               = types.MsgInviteUser
	MsgSetUserVerified          = types.MsgSetUserVerified
	MsgDepositIntoLiquidityPool = types.MsgDepositIntoLiquidityPool
	MsgAddTrustedSigner         = types.MsgAddTsp
	MsgBuyMembership            = types.MsgBuyMembership
)
