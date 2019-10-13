package memberships

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	MembershipTypeBronze = types.MembershipTypeBronze
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
	TestSetup  = keeper.SetupTestInput

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgBuyMembership = types.NewMsgBuyMembership

	// --- Tests
	TestUserAddress    = keeper.TestUserAddress
	TestMembershipType = keeper.TestMembershipType
)

type (
	Keeper           = keeper.Keeper
	Minters          = types.Minters
	MsgBuyMembership = types.MsgBuyMembership
)
